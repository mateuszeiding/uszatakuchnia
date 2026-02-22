<script lang="ts" setup>
import UCheckbox from '@/components/form/UCheckbox.vue';
import UInput from '@/components/form/UInput.vue';
import URadioButton from '@/components/form/URadioButton.vue';
</script>
<template>
    <div>
        <UInput
            name="test"
            label="test"
            type="text"
        />
        <UCheckbox
            name="test"
            label="test"
            value="1"
        />
        <UCheckbox
            name="test"
            label="test"
            value="1"
            switch
        />
        <URadioButton
            name="test"
            label="test"
            value="1"
        />
    </div>
</template>
