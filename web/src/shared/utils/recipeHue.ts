/** Stable hue (0–359) of a recipe, used by the placeholder frame of a card without a photo. */
export function recipeHue(id: number): number {
    return (((id * 47) % 360) + 360) % 360;
}
