<script lang="ts">
import { defineComponent, h } from 'vue';

// Displayed quantity: digits and unit stay in the parent's font, but the fraction glyph
// (½ ¼ ¾ ⅓ ⅔) uses the sans font: in Geist Mono it's squeezed into the width of one digit.
// Screen readers read the fraction glyph correctly, so no aria-label.
// Render function instead of a template: template formatting could add spaces around the fraction ("1 ½").
const FRAC = /([½¼¾⅓⅔])/;

export default defineComponent({
    name: 'QtyText',
    props: { text: { type: String, default: '' } },
    setup(props) {
        return () =>
            (props.text ?? '')
                .split(FRAC)
                .filter(Boolean)
                .map((t, i, all) =>
                    t.length === 1 && FRAC.test(t)
                        ? h(
                              'span',
                              {
                                  class: [
                                      'qty-frac',
                                      { 'qty-frac--joined': i > 0 && /\d$/.test(all[i - 1]!) },
                                  ],
                              },
                              t
                          )
                        : t
                );
    },
});
</script>

<style>
.qty-frac {
    font-family: var(--font-sans);
    font-size: 1.25em;
    line-height: 1;
}
.qty-frac--joined {
    margin-left: 0.04em;
}
</style>
