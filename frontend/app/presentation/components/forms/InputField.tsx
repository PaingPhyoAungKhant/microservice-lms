import React from 'react';
import { TextField, type TextFieldProps } from '@mui/material';
import type { ChangeEvent } from 'react';

type InputFieldProps = {
  label: string;
  type?: string;
  name: string;
  placeholder?: string;
  value: string;
  onChange: (e: ChangeEvent<HTMLInputElement>) => void;
  className?: string;
  error?: boolean;
  helperText?: string;
} & Omit<TextFieldProps, 'label' | 'type' | 'name' | 'value' | 'onChange'>;

const InputField: React.FC<InputFieldProps> = ({
  label,
  type = 'text',
  name,
  placeholder,
  value,
  onChange,
  className = '',
  error,
  helperText,
  ...textFieldProps
}) => {
  return (
    <TextField
      label={label}
      type={type}
      name={name}
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      className={className}
      error={error}
      helperText={helperText}
      fullWidth
      variant="outlined"
      sx={{
        mb: 2,
        '& .MuiOutlinedInput-root': {
          '& fieldset': {
            borderColor: '#D9D9D9',
          },
          '&:hover fieldset': {
            borderColor: '#1d4ed8',
          },
          '&.Mui-focused fieldset': {
            borderColor: '#1d4ed8',
            borderWidth: '1px',
          },
        },
        '& .MuiInputLabel-root': {
          fontWeight: 500,
          color: '#374151',
        },
      }}
      {...textFieldProps}
    />
  );
};

export default InputField;

