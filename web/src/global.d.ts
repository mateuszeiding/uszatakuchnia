/* eslint-disable @typescript-eslint/no-explicit-any */
export {};

type Sep = '-' | '/';
declare global {
    type NestedKeyUnion<T, Sep = '-'> = {
        [K in keyof T]: K extends string
            ? keyof T[K] extends string
                ? `${K}${Sep}${keyof T[K]}`
                : never
            : never;
    }[keyof T];

    export type PropsOf<T> = T extends keyof JSX.IntrinsicElements
        ? JSX.IntrinsicElements[T]
        : T extends new (...args: any) => { $props: infer P }
          ? P
          : T extends (...args: any) => any
            ? Parameters<T>[0] extends undefined
                ? object
                : NonNullable<Parameters<T>[0]>
            : never;
}
