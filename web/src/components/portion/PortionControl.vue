<script lang="ts" setup>
import { computed } from 'vue';

import FormControl from './FormControl.vue';
import PiecesControl from './PiecesControl.vue';

import {
    effectiveMode,
    type IPortionRecipe,
    isForm,
    isPieces,
    portionOf,
    recipeForm,
    type SizeChoice,
} from '@/shared/utils/portion';

// The size control for "sized" and "form" recipes; servings recipes keep ServingsStepper as it was.
const props = defineProps<{ recipe: IPortionRecipe; size: SizeChoice | null; big?: boolean }>();
const emit = defineEmits<{ size: [s: SizeChoice] }>();

const mode = computed(() => effectiveMode(props.recipe));
const own = computed(() => recipeForm(portionOf(props.recipe)));
</script>

<template>
    <PiecesControl
        v-if="mode === 'sized' && isPieces(size)"
        :recipe="recipe"
        :value="size"
        :big="big"
        @set="emit('size', $event)"
    />
    <FormControl
        v-else-if="mode === 'form' && isForm(size) && own"
        :value="size"
        :recipe="own"
        @set="emit('size', $event)"
    />
</template>
