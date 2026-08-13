import React from 'react';
import {
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Typography,
  Avatar,
  Divider,
  IconButton,
  Tooltip,
  useTheme
} from '@mui/material';
import {
  People as PeopleIcon,
  Assessment as ReportsIcon,
  Extension as IntegrationsIcon,
  Logout as LogoutIcon,
  Badge as BadgeIcon,
  Dashboard as DashboardIcon
} from '@mui/icons-material';
import { useLocation, useNavigate } from 'react-router-dom';
import { getCurrentUser, logout } from '../../services/auth';

const DRAWER_WIDTH = 240;

interface SideMenuProps {
  mobileOpen?: boolean;
  onMobileClose?: () => void;
}

export const SideMenu: React.FC<SideMenuProps> = ({ mobileOpen, onMobileClose }) => {
  const theme = useTheme();
  const location = useLocation();
  const navigate = useNavigate();
  const user = getCurrentUser();

  const handleLogout = () => {
    logout();
    window.location.href = '/login';
  };

  const mainItems = [
    { text: 'Karyawan', icon: <PeopleIcon />, path: '/' },
  ];

  const secondaryItems = [
    { text: 'Laporan', icon: <ReportsIcon />, path: '#', disabled: true },
    { text: 'Integrasi', icon: <IntegrationsIcon />, path: '#', disabled: true },
  ];

  const drawerContent = (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      {/* Brand Header */}
      <Box sx={{ p: 2.5, display: 'flex', alignItems: 'center', gap: 1.5 }}>
        <Box
          sx={{
            width: 32,
            height: 32,
            borderRadius: '8px',
            backgroundColor: 'primary.main',
            color: 'primary.contrastText',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <BadgeIcon sx={{ fontSize: 20 }} />
        </Box>
        <Typography variant="h6" sx={{ fontWeight: 700, letterSpacing: '-0.01em' }}>
          KaryawanApp
        </Typography>
      </Box>

      <Divider sx={{ mx: 2, opacity: 0.6 }} />

      {/* Main Items */}
      <Box sx={{ flexGrow: 1, px: 1.5, py: 2 }}>
        <Typography
          variant="caption"
          sx={{
            px: 1.5,
            pb: 1,
            display: 'block',
            fontWeight: 600,
            color: 'text.secondary',
            textTransform: 'uppercase',
            fontSize: '0.7rem',
            letterSpacing: '0.05em',
          }}
        >
          Menu Utama
        </Typography>
        <List disablePadding>
          {mainItems.map((item) => {
            const isSelected = location.pathname === item.path || (item.path === '/' && location.pathname === '/employees');
            return (
              <ListItem key={item.text} disablePadding sx={{ mb: 0.5 }}>
                <ListItemButton
                  onClick={() => {
                    navigate(item.path);
                    if (onMobileClose) onMobileClose();
                  }}
                  selected={isSelected}
                  sx={{
                    borderRadius: '8px',
                    py: 1,
                    px: 1.5,
                    color: isSelected ? 'primary.main' : 'text.primary',
                    '&.Mui-selected': {
                      backgroundColor: theme.palette.mode === 'light' ? '#F1F5F9' : '#1E293B',
                      fontWeight: 600,
                      '&:hover': {
                        backgroundColor: theme.palette.mode === 'light' ? '#E2E8F0' : '#334155',
                      },
                      '& .MuiListItemIcon-root': {
                        color: 'primary.main',
                      },
                    },
                  }}
                >
                  <ListItemIcon sx={{ minWidth: 36, color: isSelected ? 'primary.main' : 'text.secondary' }}>
                    {item.icon}
                  </ListItemIcon>
                  <ListItemText
                    primary={item.text}
                    primaryTypographyProps={{ fontSize: '0.875rem', fontWeight: isSelected ? 600 : 400 }}
                  />
                </ListItemButton>
              </ListItem>
            );
          })}
        </List>

        <Typography
          variant="caption"
          sx={{
            px: 1.5,
            pt: 2,
            pb: 1,
            display: 'block',
            fontWeight: 600,
            color: 'text.secondary',
            textTransform: 'uppercase',
            fontSize: '0.7rem',
            letterSpacing: '0.05em',
          }}
        >
          Fitur Tambahan
        </Typography>
        <List disablePadding>
          {secondaryItems.map((item) => (
            <ListItem key={item.text} disablePadding sx={{ mb: 0.5 }}>
              <ListItemButton
                disabled={item.disabled}
                sx={{
                  borderRadius: '8px',
                  py: 1,
                  px: 1.5,
                  opacity: item.disabled ? 0.5 : 1,
                }}
              >
                <ListItemIcon sx={{ minWidth: 36, color: 'text.secondary' }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText
                  primary={item.text}
                  primaryTypographyProps={{ fontSize: '0.875rem' }}
                />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </Box>

      <Divider sx={{ mx: 2, opacity: 0.6 }} />

      {/* User Profile at Bottom */}
      {user && (
        <Box sx={{ p: 2, display: 'flex', alignItems: 'center', gap: 1.5 }}>
          <Avatar
            sx={{
              width: 36,
              height: 36,
              fontSize: '0.875rem',
              fontWeight: 600,
              bgcolor: 'secondary.main',
            }}
          >
            {user.username.charAt(0).toUpperCase()}
          </Avatar>
          <Box sx={{ flexGrow: 1, minWidth: 0 }}>
            <Typography variant="body2" noWrap sx={{ fontWeight: 600 }}>
              {user.username}
            </Typography>
            <Typography variant="caption" noWrap sx={{ color: 'text.secondary', display: 'block' }}>
              {user.role}
            </Typography>
          </Box>
          <Tooltip title="Logout">
            <IconButton size="small" onClick={handleLogout} color="default">
              <LogoutIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      )}
    </Box>
  );

  return (
    <>
      {/* Desktop Permanent Drawer */}
      <Drawer
        variant="permanent"
        sx={{
          display: { xs: 'none', md: 'block' },
          width: DRAWER_WIDTH,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: DRAWER_WIDTH,
            boxSizing: 'border-box',
            borderRight: `1px solid ${theme.palette.divider}`,
            backgroundColor: theme.palette.background.paper,
          },
        }}
        open
      >
        {drawerContent}
      </Drawer>

      {/* Mobile Temporary Drawer */}
      <Drawer
        variant="temporary"
        open={mobileOpen}
        onClose={onMobileClose}
        ModalProps={{ keepMounted: true }}
        sx={{
          display: { xs: 'block', md: 'none' },
          '& .MuiDrawer-paper': {
            width: DRAWER_WIDTH,
            boxSizing: 'border-box',
            backgroundColor: theme.palette.background.paper,
          },
        }}
      >
        {drawerContent}
      </Drawer>
    </>
  );
};
