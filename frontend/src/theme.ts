import { createTheme, ThemeOptions } from '@mui/material/styles';

export const getCustomTheme = (mode: 'light' | 'dark'): ThemeOptions => ({
  palette: {
    mode,
    ...(mode === 'light'
      ? {
          primary: {
            main: '#0F172A', // Slate 900
            light: '#334155',
            dark: '#020617',
            contrastText: '#FFFFFF',
          },
          secondary: {
            main: '#6366F1', // Indigo 500
          },
          background: {
            default: '#F8FAFC', // Slate 50
            paper: '#FFFFFF',
          },
          text: {
            primary: '#0F172A',
            secondary: '#64748B',
          },
          divider: '#E2E8F0',
        }
      : {
          primary: {
            main: '#F8FAFC',
            light: '#FFFFFF',
            dark: '#CBD5E1',
            contrastText: '#0F172A',
          },
          secondary: {
            main: '#818CF8',
          },
          background: {
            default: '#0B0F19',
            paper: '#111827',
          },
          text: {
            primary: '#F9FAFB',
            secondary: '#9CA3AF',
          },
          divider: '#1F2937',
        }),
  },
  typography: {
    fontFamily: [
      'Inter',
      '-apple-system',
      'BlinkMacSystemFont',
      '"Segoe UI"',
      'Roboto',
      'sans-serif',
    ].join(','),
    h4: {
      fontWeight: 700,
      fontSize: '1.75rem',
      letterSpacing: '-0.02em',
    },
    h5: {
      fontWeight: 600,
      fontSize: '1.25rem',
      letterSpacing: '-0.01em',
    },
    h6: {
      fontWeight: 600,
      fontSize: '1rem',
    },
    subtitle2: {
      fontWeight: 500,
      fontSize: '0.875rem',
    },
    button: {
      textTransform: 'none',
      fontWeight: 600,
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          boxShadow: 'none',
          '&:hover': {
            boxShadow: 'none',
          },
        },
        containedPrimary: {
          backgroundColor: mode === 'light' ? '#0F172A' : '#F8FAFC',
          color: mode === 'light' ? '#FFFFFF' : '#0F172A',
          '&:hover': {
            backgroundColor: mode === 'light' ? '#1E293B' : '#E2E8F0',
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
          boxShadow: mode === 'light' 
            ? '0px 1px 3px rgba(15, 23, 42, 0.03), 0px 1px 2px rgba(15, 23, 42, 0.06)'
            : '0px 1px 3px rgba(0, 0, 0, 0.3)',
          border: `1px solid ${mode === 'light' ? '#E2E8F0' : '#1F2937'}`,
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          border: `1px solid ${mode === 'light' ? '#E2E8F0' : '#1F2937'}`,
          boxShadow: 'none',
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          fontWeight: 500,
          borderRadius: 6,
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        head: {
          fontWeight: 600,
          backgroundColor: mode === 'light' ? '#F8FAFC' : '#111827',
          color: mode === 'light' ? '#475569' : '#9CA3AF',
          borderBottom: `1px solid ${mode === 'light' ? '#E2E8F0' : '#1F2937'}`,
        },
        root: {
          borderBottom: `1px solid ${mode === 'light' ? '#F1F5F9' : '#1F2937'}`,
          fontSize: '0.875rem',
        },
      },
    },
  },
});

const theme = createTheme(getCustomTheme('light'));

export { theme };
