// Auto-generated types from colors.config.mjs
export const palette: {
  "gray": {
    "ghost-white": string;
    "cool-gray": string;
    "paynes-gray": string;
    "outer-space": string;
    "davys-gray": string;
    "eerie-black": string;
  };
  "blue": {
    "deep-sky-blue": string;
    "picton-blue": string;
    "steel-blue": string;
    "celtic-blue": string;
    "cerulean": string;
    "indigo-dye": string;
    "berkeley-blue": string;
    "delft-blue": string;
  };
  "red": {
    "bright-pink": string;
    "old-rose": string;
    "redwood": string;
    "china-rose": string;
    "imperial-red": string;
    "chili-red": string;
    "crimson": string;
    "amaranth-purple": string;
    "quinacridone-magenta": string;
    "falu-red": string;
  };
  "green": {
    "spring-green": string;
    "cambridge-blue": string;
    "celadon": string;
    "tiffany-blue": string;
    "mint": string;
    "stale-gray": string;
    "fern-green": string;
    "moss-green": string;
    "myrtle-green": string;
    "persian-green": string;
    "dartmouth-green": string;
    "phlato-green": string;
    "drab-dark-brown": string;
  };
  "purple": {
    "mardi-gras": string;
    "royal-purple": string;
    "finn": string;
    "dark-purple": string;
  };
  "violet": {
    "plum": string;
    "chinese-violet": string;
    "french-violet": string;
    "ultra-violet": string;
    "russian-violet": string;
  };
  "yellow": {
    "vanilla": string;
    "maize": string;
    "old-gold": string;
    "selective-yellow": string;
    "xanthous": string;
  };
  "brown": {
    "lavender-blush": string;
    "sandy-brown": string;
    "wenge": string;
    "taupe": string;
    "licorice": string;
  };
  "orange": {
    "ut-orange": string;
    "cocoa-brown": string;
    "burnt-orange": string;
  };
};

export type Palette = typeof palette;
export type PaletteGroup = keyof Palette;
export type PaletteKey<G extends PaletteGroup = PaletteGroup> = keyof Palette[G];