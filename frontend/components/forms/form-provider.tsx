"use client";

import { FormProvider, UseFormReturn } from 'react-hook-form';
import { PropsWithChildren } from 'react';

type FormWrapperProps = PropsWithChildren<{ methods: UseFormReturn<any>; onSubmit: (values: any) => void }>;

export function VFormProvider({ methods, onSubmit, children }: FormWrapperProps) {
  return (
    <FormProvider {...methods}>
      <form onSubmit={methods.handleSubmit(onSubmit)} style={{ display: 'grid', gap: 16 }}>
        {children}
      </form>
    </FormProvider>
  );
}
