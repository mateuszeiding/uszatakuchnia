<script lang="ts" setup>
import UCheckboxGroup from '@/components/form/UCheckboxGroup.vue';
import UInputText from '@/components/form/UInputText.vue';
import URadioGroup from '@/components/form/URadioGroup.vue';
import USwitch from '@/components/form/USwitch.vue';
import UTextarea from '@/components/form/UTextarea.vue';
</script>
<template>
    <div>
        <UInputText
            name="test"
            label="test"
        />
        <UTextarea
            name="test"
            label="test"
        />
        <UCheckboxGroup
            name="test"
            legend="Group test"
            :options="[
                {
                    label: 'one',
                    value: 'one',
                },
                {
                    label: 'two',
                    value: 'two',
                },
            ]"
        />
        <USwitch
            name="test"
            label="test"
        />
        <URadioGroup
            name="test"
            legend="Group test"
            :options="[
                {
                    label: 'one',
                    value: 'one',
                },
                {
                    label: 'two',
                    value: 'two',
                },
            ]"
        />
    </div>
</template>
