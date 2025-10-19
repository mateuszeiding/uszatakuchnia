import { writeFile, mkdir } from "node:fs/promises";
import { resolve, dirname, sep } from "node:path";
import { pathToFileURL } from "node:url";

const rawArgs = process.argv.slice(2);
const args = {};
for (let i = 0; i < rawArgs.length; i++) {
    const a = rawArgs[i];
    if (a.startsWith("--")) {
        const key = a.slice(2);
        const val = rawArgs[i + 1] && !rawArgs[i + 1].startsWith("--") ? rawArgs[++i] : true;
        args[key] = val;
    }
}

const CWD = process.cwd();
const PALETTE_PATH = resolve(CWD, args.palette ?? "./colors.config.mjs");
const OUT_SCSS = resolve(CWD, args["out-scss"] ?? "./config/_colors.scss");
const OUT_MD = resolve(CWD, args["out-md"] ?? "./colors.md");

const isHex = (v) => /^#([0-9a-f]{6}|[0-9a-f]{3})$/i.test(v);
const normHex = (h) =>
    h.length === 4 ? "#" + h[1] + h[1] + h[2] + h[2] + h[3] + h[3] : h.toLowerCase();

function sortKeys(keys) {
    const allNum = keys.every((k) => /^\d+$/.test(k));
    if (allNum) return [...keys].sort((a, b) => Number(b) - Number(a));
    return [...keys].sort();
}

async function ensureDirFor(filePath) {
    await mkdir(dirname(filePath), { recursive: true });
}

function relFromCwd(abs) {
    const prefix = CWD.endsWith(sep) ? CWD : CWD + sep;
    return abs.startsWith(prefix) ? abs.slice(prefix.length) : abs;
}

// ---------- generatory ----------
function toScss(palette) {
    const groups = Object.keys(palette);
    const out = [];

    out.push(
        `@use 'sass:map';
@use 'sass:color';
@use 'sass:meta';

// ⚠️ GENERATED from ${relFromCwd(PALETTE_PATH)} — do not edit.
`
    );

    for (const g of groups) {
        out.push(`/* --- ${g.toUpperCase()} --- */`);
        const keys = sortKeys(Object.keys(palette[g]));
        for (const k of keys) {
            const hex = normHex(palette[g][k]);
            if (!isHex(hex)) throw new Error(`Invalid HEX: ${g}.${k}=${palette[g][k]}`);
            out.push(`$${g}-${k}: ${hex};`);
        }
        out.push("");
    }

    for (const g of groups) {
        const keys = sortKeys(Object.keys(palette[g]));
        const entries = keys.map((k) => `  ${g}-${k}: $${g}-${k},`).join("\n");
        out.push(`$${g}s: (\n${entries}\n);\n`);
    }

    out.push(
        `@function hex-to-rgb($c) {
  @return rgb(
    color.channel($c, 'red', $space: rgb),
    color.channel($c, 'green', $space: rgb),
    color.channel($c, 'blue', $space: rgb)
  );
}

@function map-merge-all($maps...) {
  $out: ();
  @each $m in $maps {
    @if $m != null and meta.type-of($m) == 'map' {
      $out: map.merge($out, $m);
    }
  }
  @return $out;
}
`
    );

    const mapList = groups.map((g) => `$${g}s`).join(", ");
    out.push(`$colors: map-merge-all(${mapList});\n`);

    out.push(`:root {`);
    for (const g of groups) {
        const keys = sortKeys(Object.keys(palette[g]));
        for (const k of keys) {
            out.push(`  --color-${g}-${k}: #{$${g}-${k}};`);
        }
    }
    out.push(`}\n`);

    out.push(
        `@each $name, $value in $colors {
  .tx-#{$name} { color: var(--color-#{$name}); }
  .bg-#{$name} { background-color: var(--color-#{$name}); }
}
`
    );

    return out.join("\n");
}

function unslug(name) {
    return name.replace(/-/g, " ");
}

function toMarkdown(palette) {
    const groups = Object.keys(palette);
    const out = [];

    out.push(
        `# colors 🎨

## Colors from [https://coolors.co](https://coolors.co)
`
    );

    for (const g of groups) {
        out.push(`### ${g.toUpperCase()} ![${g}](https://placehold.co/20x20/${g}/${g}.png)\n`);
        const keys = sortKeys(Object.keys(palette[g]));
        for (const k of keys) {
            const hex = normHex(palette[g][k]).toUpperCase();
            const displayName = unslug(k); // zamiast gray-outer-space → outer space
            const label = encodeURIComponent(displayName);
            const hexEsc = hex.replace("#", "%23");
            out.push(
                `![${displayName}](https://img.shields.io/badge/${label}-%23333?style=for-the-badge&label=${hexEsc}&labelColor=${hexEsc})`
            );
        }
        out.push(
            "\n"
        );

    }

    out.push(`---`);
    return out.join("\n");
}

// ---------- main ----------
(async function main() {
    const mod = await import(pathToFileURL(PALETTE_PATH).href);
    const palette = mod.palette;
    if (!palette || typeof palette !== "object") {
        throw new Error(`Nie znaleziono exportu "palette" w ${PALETTE_PATH}`);
    }

    for (const [g, shades] of Object.entries(palette)) {
        for (const [k, hex] of Object.entries(shades)) {
            if (!isHex(hex)) throw new Error(`Invalid HEX in palette: ${g}.${k} = ${hex}`);
        }
    }

    await ensureDirFor(OUT_SCSS);
    await ensureDirFor(OUT_MD);

    await writeFile(OUT_SCSS, toScss(palette), "utf8");
    await writeFile(OUT_MD, toMarkdown(palette), "utf8");

    console.log("✅ Generated:");
    console.log(" -", relFromCwd(OUT_SCSS));
    console.log(" -", relFromCwd(OUT_MD));
})().catch((err) => {
    console.error("❌ genColors failed:", err);
    process.exit(1);
});


const OUT = resolve("./colors.config.d.ts");

const mod = await import(pathToFileURL(PALETTE_PATH).href);
const palette = mod.palette;

function genTypes(obj) {
    const lines = [];
    lines.push("// Auto-generated types from colors.config.mjs");
    lines.push("export const palette: {");
    for (const [group, shades] of Object.entries(obj)) {
        lines.push(`  ${JSON.stringify(group)}: {`);
        for (const [key] of Object.entries(shades)) {
            lines.push(`    ${JSON.stringify(key)}: string;`);
        }
        lines.push("  };");
    }
    lines.push("};");
    lines.push("");
    lines.push("export type Palette = typeof palette;");
    lines.push("export type PaletteGroup = keyof Palette;");
    lines.push("export type PaletteKey<G extends PaletteGroup = PaletteGroup> = keyof Palette[G];");
    return lines.join("\n");
}

await writeFile(OUT, genTypes(palette), "utf8");
console.log("✅ Generated types:", OUT);