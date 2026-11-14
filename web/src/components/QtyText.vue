<script lang="ts">
import { defineComponent, h } from 'vue';

// Ilość do wyświetlenia: cyfry i jednostka zostają w foncie rodzica (Geist Mono), a znak ułamka
// (½ ¼ ¾ ⅓ ⅔) idzie w Geist sans 1,25em — w Geist Mono glif jest wciśnięty w szerokość jednej cyfry.
// Czytnik ekranu czyta znak ułamka poprawnie, więc bez aria-label.
// Render function zamiast szablonu: formatowanie szablonu nie może dokleić spacji wokół ułamka („1 ½”).
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
