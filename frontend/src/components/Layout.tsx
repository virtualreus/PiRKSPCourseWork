import { Link, Outlet } from 'react-router-dom';

import { useAuth } from '../context/AuthContext';

export function Layout() {
  const { user, logout } = useAuth();

  return (
    <div className="app">
      <header className="header">
        <Link to="/" className="logo">
          Hackathon<span>Hub</span>
        </Link>
        <nav className="nav">
          {user ? (
            <>
              <Link to="/profile">Профиль</Link>
              {user.platform_role === 'organizer' && (
                <span className="nav-badge">организатор</span>
              )}
              <button type="button" className="btn-ghost" onClick={logout}>
                Выйти
              </button>
            </>
          ) : (
            <>
              <Link to="/login">Вход</Link>
              <Link to="/register" className="btn-primary btn-sm">
                Регистрация
              </Link>
            </>
          )}
        </nav>
      </header>
      <main className="main">
        <Outlet />
      </main>
    </div>
  );
}
