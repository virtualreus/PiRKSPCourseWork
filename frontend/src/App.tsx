import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';

import { GuestRoute } from './components/GuestRoute';
import { Layout } from './components/Layout';
import { OrganizerRoute } from './components/OrganizerRoute';
import { ProtectedRoute } from './components/ProtectedRoute';
import { AuthProvider } from './context/AuthContext';
import { HackathonDetailPage } from './pages/HackathonDetailPage';
import { HomePage } from './pages/HomePage';
import { LoginPage } from './pages/LoginPage';
import { ProfilePage } from './pages/ProfilePage';
import { RegisterPage } from './pages/RegisterPage';
import { OrganizerHackathonEditPage } from './pages/organizer/OrganizerHackathonEditPage';
import { OrganizerHackathonNewPage } from './pages/organizer/OrganizerHackathonNewPage';
import { OrganizerHackathonsPage } from './pages/organizer/OrganizerHackathonsPage';

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route element={<Layout />}>
            <Route index element={<HomePage />} />
            <Route path="hackathons/:id" element={<HackathonDetailPage />} />
            <Route
              path="login"
              element={
                <GuestRoute>
                  <LoginPage />
                </GuestRoute>
              }
            />
            <Route
              path="register"
              element={
                <GuestRoute>
                  <RegisterPage />
                </GuestRoute>
              }
            />
            <Route
              path="profile"
              element={
                <ProtectedRoute>
                  <ProfilePage />
                </ProtectedRoute>
              }
            />
            <Route
              path="organizer/hackathons"
              element={
                <OrganizerRoute>
                  <OrganizerHackathonsPage />
                </OrganizerRoute>
              }
            />
            <Route
              path="organizer/hackathons/new"
              element={
                <OrganizerRoute>
                  <OrganizerHackathonNewPage />
                </OrganizerRoute>
              }
            />
            <Route
              path="organizer/hackathons/:id"
              element={
                <OrganizerRoute>
                  <OrganizerHackathonEditPage />
                </OrganizerRoute>
              }
            />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}
