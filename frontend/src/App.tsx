import React, { useState, useEffect } from 'react';
import { 
  BrowserRouter as Router, 
  Routes, 
  Route, 
  Navigate, 
  useLocation, 
  Link as RouterLink,
  Outlet
} from 'react-router-dom';
import { 
  Container, 
  AppBar, 
  Toolbar, 
  Typography, 
  Box, 
  Button,
  CircularProgress
} from '@mui/material';
import EmployeeList from './components/EmployeeList';
import EmployeeForm from './components/EmployeeForm';
import Login from './pages/Login';
import ProtectedRoute from './components/ProtectedRoute';
import { getCurrentUser, logout } from './services/auth';

const Navigation = () => {
  const location = useLocation();
  const [user, setUser] = useState(getCurrentUser());
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setUser(getCurrentUser());
  }, [location]);

  const handleLogout = () => {
    try {
      setLoading(true);
      logout();
      window.location.href = '/login';
    } catch (error) {
      console.error('Logout failed:', error);
      setLoading(false);
    }
  };

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          <RouterLink to="/" style={{ color: 'white', textDecoration: 'none' }}>
            Employee Management System
          </RouterLink>
        </Typography>
        {user ? (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <Typography variant="subtitle1">
              Welcome, {user.username} ({user.role})
            </Typography>
            <Button 
              variant="outlined" 
              color="inherit" 
              onClick={handleLogout}
              disabled={loading}
              startIcon={loading ? <CircularProgress size={20} color="inherit" /> : null}
              sx={{ color: 'white', borderColor: 'white' }}
            >
              {loading ? 'Logging out...' : 'Logout'}
            </Button>
          </Box>
        ) : (
          <Button 
            color="inherit" 
            component={RouterLink} 
            to="/login"
            state={{ from: location.pathname }}
          >
            Login
          </Button>
        )}
      </Toolbar>
    </AppBar>
  );
};

// Protected layout component
const ProtectedLayout: React.FC<{ children?: React.ReactNode }> = ({ children }) => (
  <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
    <Navigation />
    <Container component="main" sx={{ mt: 4, mb: 4, flex: 1 }}>
      {children || <Outlet />}
    </Container>
  </Box>
);

function App() {
  return (
    <Router>
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={<Login />} />
        
        {/* Protected routes */}
        <Route element={
          <ProtectedRoute>
            <ProtectedLayout />
          </ProtectedRoute>
        }>
          <Route path="/" element={<EmployeeList />} />
          <Route path="/employees" element={<EmployeeList />} />
          <Route path="/employees/new" element={<EmployeeForm />} />
          <Route path="/employees/:id/edit" element={<EmployeeForm />} />
        </Route>
        
        {/* Catch all other routes */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Router>
  );
}

export default App;
