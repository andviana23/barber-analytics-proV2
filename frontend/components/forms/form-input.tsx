"use client";

import { Controller, useFormContext } from 'react-hook-form';
import { VTextField } from '../ui/val-textfield';
import { TextFieldProps } from '@mui/material';

type FormInputProps = {
  name: string;
  label: string;
  helperText?: string;
} & Omit<TextFieldProps, 'name' | 'label'>;

export function FormInput({ name, label, helperText, ...rest }: FormInputProps) {
  const { control } = useFormContext();

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => (
        <VTextField
          {...rest}
          {...field}
          label={label}
          error={!!fieldState.error}
          helperText={fieldState.error?.message || helperText}
        />
      )}
    />
  );
}
