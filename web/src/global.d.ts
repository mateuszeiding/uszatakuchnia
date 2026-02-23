/* eslint-disable @typescript-eslint/no-explicit-any */
export {};

declare global {
    type Sep = '-' | '/';
    type NestedKeyUnion<T, Sep extends string = '/'> = {
        [K in Extract<keyof T, string>]: T[K] extends object
            ? `${K}${Sep}${KeyToString<Exclude<keyof T[K], symbol>>}`
            : never;
    }[Extract<keyof T, string>];

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
