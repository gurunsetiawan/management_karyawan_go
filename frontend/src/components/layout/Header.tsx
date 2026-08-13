import React from 'react';
import {
  Box,
  IconButton,
  Typography,
  Breadcrumbs,
  Tooltip,
  useTheme
} from '@mui/material';
import {
  Menu as MenuIcon,
  Brightness4 as DarkModeIcon,
  Brightness7 as LightModeIcon,
  NavigateNext as NavigateNextIcon
} from '@mui/icons-material';
import { useLocation } from 'react-router-dom';

interface HeaderProps {
  onMobileMenuToggle: () => void;
  mode: 'light' | 'dark';
  onToggleMode: () => void;
}

export const Header: React.FC<HeaderProps> = ({ onMobileMenuToggle, mode, onToggleMode }) => {
  const theme = useTheme();
  const location = useLocation();

  const getBreadcrumbs = () => {
    if (location.pathname.includes('/new')) {
      return ['Karyawan', 'Tambah Karyawan'];
    }
    if (location.pathname.includes('/edit')) {
      return ['Karyawan', 'Edit Karyawan'];
    }
    return ['Employees'];
  };

  const crumbs = getBreadcrumbs();

  return (
    <Box
      component="header"
      sx={{
        py: 1.5,
        px: { xs: 2, md: 4 },
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        backgroundColor: theme.palette.background.default,
        borderBottom: `1px solid ${theme.palette.divider}`,
        position: 'sticky',
        top: 0,
        zIndex: 1100,
      }}
    >
      {/* Left: Mobile Menu & Breadcrumbs */}
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <IconButton
          color="inherit"
          aria-label="open drawer"
          edge="start"
          onClick={onMobileMenuToggle}
          sx={{ display: { md: 'none' } }}
        >
          <MenuIcon />
        </IconButton>

        <Breadcrumbs
          separator={<NavigateNextIcon fontSize="small" sx={{ color: 'text.secondary' }} />}
          aria-label="breadcrumb"
        >
          {crumbs.map((crumb, idx) => (
            <Typography
              key={crumb}
              variant="body2"
              sx={{
                color: idx === crumbs.length - 1 ? 'text.primary' : 'text.secondary',
                fontWeight: idx === crumbs.length - 1 ? 600 : 400,
              }}
            >
              {crumb}
            </Typography>
          ))}
        </Breadcrumbs>
      </Box>

      {/* Right: Theme Switcher & Options */}
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Tooltip title={mode === 'light' ? 'Mode Gelap' : 'Mode Terang'}>
          <IconButton size="small" onClick={onToggleMode} color="inherit">
            {mode === 'light' ? <DarkModeIcon fontSize="small" /> : <LightModeIcon fontSize="small" />}
          </IconButton>
        </Tooltip>
      </Box>
    </Box>
  );
};
