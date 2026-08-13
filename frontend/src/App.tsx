import React, { useState, useMemo } from 'react';
import { 
  BrowserRouter as Router, 
  Routes, 
  Route, 
  Navigate 
} from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline } from '@mui/material';
import { getCustomTheme } from './theme';
import EmployeeList from './components/EmployeeList';
import EmployeeForm from './components/EmployeeForm';
import Login from './pages/Login';
import ProtectedRoute from './components/ProtectedRoute';
import { DashboardLayout } from './components/layout/DashboardLayout';

function App() {
  const [mode, setMode] = useState<'light' | 'dark'>('light');

  const theme = useMemo(() => createTheme(getCustomTheme(mode)), [mode]);

  const toggleMode = () => {
    setMode((prevMode) => (prevMode === 'light' ? 'dark' : 'light'));
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Routes>
          {/* Public routes */}
          <Route path="/login" element={<Login />} />
          
          {/* Protected routes wrapped in DashboardLayout */}
          <Route
            element={
              <ProtectedRoute>
                <DashboardLayout mode={mode} onToggleMode={toggleMode} />
              </ProtectedRoute>
            }
          >
            <Route path="/" element={<EmployeeList />} />
            <Route path="/employees" element={<EmployeeList />} />
            <Route path="/employees/new" element={<EmployeeForm />} />
            <Route path="/employees/:id/edit" element={<EmployeeForm />} />
          </Route>
          
          {/* Catch all other routes */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App;
