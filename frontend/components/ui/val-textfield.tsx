"use client";

import { TextField, TextFieldProps } from '@mui/material';

export function VTextField(props: TextFieldProps) {
  return (
    <TextField
      {...props}
      fullWidth
      sx={{
        '& .MuiOutlinedInput-root': {
          borderRadius: '10px',
          backgroundColor: 'var(--valtaris-surface-subtle)',
        },
        '& .MuiOutlinedInput-notchedOutline': {
          borderColor: 'var(--valtaris-border)',
        },
        '&:hover .MuiOutlinedInput-notchedOutline': {
          borderColor: 'var(--valtaris-primary)',
        },
        '& .Mui-error .MuiOutlinedInput-notchedOutline': {
          borderColor: 'var(--valtaris-danger)',
        },
        ...(props.sx || {}),
      }}
    />
  );
}
