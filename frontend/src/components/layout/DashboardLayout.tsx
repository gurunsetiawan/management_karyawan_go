import React, { useState } from 'react';
import { Box } from '@mui/material';
import { SideMenu } from './SideMenu';
import { Header } from './Header';
import { Outlet } from 'react-router-dom';

const DRAWER_WIDTH = 240;

interface DashboardLayoutProps {
  children?: React.ReactNode;
  mode: 'light' | 'dark';
  onToggleMode: () => void;
}

export const DashboardLayout: React.FC<DashboardLayoutProps> = ({
  children,
  mode,
  onToggleMode,
}) => {
  const [mobileOpen, setMobileOpen] = useState(false);

  const handleMobileMenuToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh', backgroundColor: 'background.default' }}>
      {/* Sidebar */}
      <SideMenu mobileOpen={mobileOpen} onMobileClose={() => setMobileOpen(false)} />

      {/* Main Content Body */}
      <Box
        sx={{
          flexGrow: 1,
          display: 'flex',
          flexDirection: 'column',
          minWidth: 0,
          width: { md: `calc(100% - ${DRAWER_WIDTH}px)` },
        }}
      >
        <Header
          onMobileMenuToggle={handleMobileMenuToggle}
          mode={mode}
          onToggleMode={onToggleMode}
        />

        <Box
          component="main"
          sx={{
            flexGrow: 1,
            p: { xs: 2, sm: 3, md: 4 },
            maxWidth: 1400,
            width: '100%',
            mx: 'auto',
          }}
        >
          {children || <Outlet />}
        </Box>
      </Box>
    </Box>
  );
};
