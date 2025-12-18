import React from 'react';
import { Button as MuiButton, type ButtonProps as MuiButtonProps } from '@mui/material';
import { styled } from '@mui/material/styles';

type ButtonVariant = 'fill' | 'outline';
type ButtonSize = 'sm' | 'md' | 'lg';
type ButtonColor = 'primary' | 'secondary';
type ButtonRounded = 'none' | 'sm' | 'md' | 'lg' | 'full';

interface CustomButtonProps extends Omit<MuiButtonProps, 'variant' | 'size' | 'color'> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  color?: ButtonColor;
  rounded?: ButtonRounded;
}

const StyledButton = styled(MuiButton, {
  shouldForwardProp: (prop) => prop !== 'customVariant' && prop !== 'customColor' && prop !== 'rounded',
})<{ customVariant?: ButtonVariant; customColor?: ButtonColor; rounded?: ButtonRounded }>(
  ({ theme, customVariant, customColor, rounded }) => ({
    textTransform: 'none',
    fontWeight: 500,
    transition: 'all 0.2s',
    '&:hover': {
      transform: 'scale(1.05)',
    },
    ...(customVariant === 'fill' && {
      backgroundColor:
        customColor === 'primary' ? theme.palette.primary.main : theme.palette.secondary.main,
      color: customColor === 'primary' ? theme.palette.primary.contrastText : theme.palette.secondary.contrastText,
      '&:hover': {
        backgroundColor:
          customColor === 'primary' ? theme.palette.primary.dark : theme.palette.secondary.dark,
        transform: 'scale(1.05)',
      },
    }),
    ...(customVariant === 'outline' && {
      borderWidth: '2px',
      borderStyle: 'solid',
      borderColor: customColor === 'primary' ? theme.palette.primary.main : theme.palette.secondary.main,
      color: customColor === 'primary' ? theme.palette.primary.main : theme.palette.secondary.main,
      backgroundColor: 'transparent',
      '&:hover': {
        borderColor: customColor === 'primary' ? theme.palette.primary.dark : theme.palette.secondary.dark,
        backgroundColor: 'transparent',
        transform: 'scale(1.05)',
      },
    }),
    ...(rounded === 'none' && { borderRadius: 0 }),
    ...(rounded === 'sm' && { borderRadius: '0.125rem' }),
    ...(rounded === 'md' && { borderRadius: '0.375rem' }),
    ...(rounded === 'lg' && { borderRadius: '0.5rem' }),
    ...(rounded === 'full' && { borderRadius: '9999px' }),
  })
);

const sizeMap: Record<ButtonSize, { padding: string; fontSize: string }> = {
  sm: { padding: '1rem', fontSize: '0.875rem' },
  md: { padding: '1.5rem', fontSize: '1rem' },
  lg: { padding: '2rem 1rem 0.75rem', fontSize: '1.125rem' },
};

const Button: React.FC<CustomButtonProps> = ({
  children,
  variant = 'fill',
  color = 'primary',
  size = 'md',
  rounded = 'md',
  className = '',
  fullWidth,
  sx,
  ...props
}) => {
  const sizeStyles = sizeMap[size];

  return (
    <StyledButton
      customVariant={variant}
      customColor={color}
      rounded={rounded}
      className={className}
      fullWidth={fullWidth}
      sx={{
        ...sizeStyles,
        ...sx,
      }}
      {...props}
    >
      {children}
    </StyledButton>
  );
};

export default Button;

