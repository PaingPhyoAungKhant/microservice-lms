import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router';
import {
  Avatar,
  Menu,
  MenuItem,
  IconButton,
  Divider,
  Typography,
  Box,
} from '@mui/material';
import DashboardIcon from '@mui/icons-material/Dashboard';
import LogoutIcon from '@mui/icons-material/Logout';
import Button from '../common/Button';
import { useAuth } from '../../hooks/useAuth';
import { ROUTES } from '../../../shared/constants/routes';
import { ROLE_LABELS } from '../../../shared/constants/roles';

function HamburgerMenu() {
  const [isOpen, setIsOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const toggleMenu = () => setIsOpen(!isOpen);

  const handleUserMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleUserMenuClose = () => {
    setAnchorEl(null);
  };

  const handleGoToDashboard = () => {
    handleUserMenuClose();
    if (user?.role === 'admin') {
      navigate(ROUTES.ADMIN_USERS);
    } else if (user?.role === 'instructor') {
      navigate(ROUTES.INSTRUCTOR_DASHBOARD);
    } else if (user?.role === 'student') {
      navigate(ROUTES.STUDENT_DASHBOARD);
    }
  };

  const handleLogout = () => {
    handleUserMenuClose();
    logout();
  };

  return (
    <div className="relative w-full">
      {/* Toggle Menu */}
      <div
        className="absolute right-0 z-100 flex w-6 flex-col gap-1.5 md:hidden"
        onClick={toggleMenu}
      >
        <span
          className={`h-0.5 w-full transform bg-black transition-all duration-300 ${isOpen ? 'translate-y-2 rotate-45' : ''}`}
        ></span>
        <span
          className={`h-0.5 w-full bg-black transition-all duration-300 ${isOpen ? 'opacity-0' : 'opacity-100'}`}
        ></span>
        <span
          className={`h-0.5 w-full bg-black transition-all duration-300 ${isOpen ? '-translate-y-2 -rotate-45' : ''}`}
        ></span>
      </div>
      {/* Desktop Navigation */}
      <div className="ml-6 hidden flex-row items-center justify-between space-x-6 text-xl md:flex">
        <div className="flex flex-row">
          <Link
            to={ROUTES.HOME}
            className="hover:text-primary-2 text-text-primary block w-full rounded-lg px-4 py-2 transition-all duration-200 hover:scale-105 focus:outline-none"
            onClick={toggleMenu}
            viewTransition
          >
            Home
          </Link>

          <Link
            to={ROUTES.COURSES}
            className="hover:text-primary-2 text-text-primary block w-full rounded-lg px-4 py-2 transition-all duration-200 hover:scale-105 focus:outline-none"
            onClick={toggleMenu}
            viewTransition
          >
            Courses
          </Link>
        </div>

        <div className="flex flex-row items-center">
          {user ? (
            <>
              <IconButton
                onClick={handleUserMenuOpen}
                size="small"
                sx={{ ml: 1 }}
                aria-label="user menu"
              >
                <Avatar sx={{ width: 32, height: 32, bgcolor: 'primary.main' }}>
                  {user.username?.[0]?.toUpperCase() || 'U'}
                </Avatar>
              </IconButton>
              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleUserMenuClose}
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'right',
                }}
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
              >
                <MenuItem disabled>
                  <Box>
                    <Typography variant="body2" fontWeight={600}>
                      {user.username}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {user.email}
                    </Typography>
                    <Typography variant="caption" display="block" color="text.secondary">
                      {user.role ? ROLE_LABELS[user.role] : ''}
                    </Typography>
                  </Box>
                </MenuItem>
                <Divider />
                <MenuItem onClick={handleGoToDashboard}>
                  <DashboardIcon fontSize="small" sx={{ mr: 1 }} />
                  <Typography variant="body2">Go to Dashboard</Typography>
                </MenuItem>
                <Divider />
                <MenuItem onClick={handleLogout}>
                  <LogoutIcon fontSize="small" sx={{ mr: 1 }} />
                  <Typography variant="body2">Logout</Typography>
                </MenuItem>
              </Menu>
            </>
          ) : (
            <>
              <Link
                to={ROUTES.LOGIN}
                className="hover:text-primary-2 text-text-primary rounded-lg px-4 py-2 transition-all duration-200 hover:scale-105 focus:outline-none"
                onClick={toggleMenu}
                viewTransition
              >
                Login
              </Link>
              <Link
                to={ROUTES.REGISTER}
                className="hover:text-primary-2 text-text-primary cursor-pointer rounded-lg px-4 py-2 transition-all duration-200 hover:scale-105 focus:outline-none"
                onClick={toggleMenu}
                viewTransition
              >
                <Button size="sm" variant="outline" color="secondary">
                  Join ASTO Academy
                </Button>
              </Link>
            </>
          )}
        </div>
      </div>

      {/* Mobile Navigation Menu Overlay */}

      <div
        className={`fixed inset-0 z-50 transform transition-opacity duration-300 md:hidden ${isOpen ? 'opacity-100' : 'pointer-events-none opacity-0'}`}
        onClick={toggleMenu}
      >
        <div
          className={`bg-primary-3 absolute top-0 right-0 h-full w-3/4 transform shadow-xl transition-transform duration-300 ease-in-out ${isOpen ? 'translate-x-0' : 'translate-x-full'}`}
          onClick={(e) => e.stopPropagation()}
        >
          <div className="mt-26 flex flex-col gap-2 px-6">
            <Link
              to={ROUTES.HOME}
              className="hover:text-primary-1 block w-full rounded-lg px-4 py-2 transition-all duration-200 hover:scale-105 focus:outline-none"
              onClick={toggleMenu}
              viewTransition
            >
              Home
            </Link>

            <Link
              to={ROUTES.COURSES}
              className="hover:text-primary-1 block w-full rounded-lg px-4 py-2 transition-all duration-200 hover:scale-105 focus:outline-none"
              onClick={toggleMenu}
              viewTransition
            >
              Courses
            </Link>
            {user ? (
              <>
                <div className="px-4 py-2 border-b border-gray-200 mb-2">
                  <div className="font-semibold text-base">{user.username}</div>
                  <div className="text-sm text-gray-500">{user.email}</div>
                  <div className="text-xs text-gray-500">{user.role ? ROLE_LABELS[user.role] : ''}</div>
                </div>
                <Link
                  to={
                    user.role === 'admin'
                      ? ROUTES.ADMIN_USERS
                      : user.role === 'instructor'
                        ? ROUTES.INSTRUCTOR_DASHBOARD
                        : ROUTES.STUDENT_DASHBOARD
                  }
                  className="hover:text-primary-1 block w-full rounded-lg px-4 py-2 transition-all duration-200 hover:scale-105 focus:outline-none"
                  onClick={toggleMenu}
                  viewTransition
                >
                  Go to Dashboard
                </Link>
                <button
                  onClick={() => {
                    toggleMenu();
                    logout();
                  }}
                  className="hover:text-primary-1 block w-full rounded-lg px-4 py-2 text-left transition-all duration-200 hover:scale-105 focus:outline-none"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link
                  to={ROUTES.LOGIN}
                  className="hover:text-primary-1 block w-full rounded-lg px-4 py-2 transition-all duration-200 hover:scale-105 focus:outline-none"
                  onClick={toggleMenu}
                  viewTransition
                >
                  Login
                </Link>
                <Link
                  to={ROUTES.REGISTER}
                  className="hover:text-primary-1 block w-full rounded-lg px-4 py-2 transition-all duration-200 hover:scale-105 focus:outline-none"
                  onClick={toggleMenu}
                  viewTransition
                >
                  Register
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default HamburgerMenu;

