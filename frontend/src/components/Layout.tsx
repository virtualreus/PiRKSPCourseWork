import { useState } from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';

import { useAuth } from '../context/AuthContext';
import { SceneBackground } from './SceneBackground';

export function Layout() {
  const { user, logout } = useAuth();
  const [menuOpen, setMenuOpen] = useState(false);
  const location = useLocation();

  return (
    <div className="app">
      <SceneBackground />
      <header className="header glass">
        <div className="header-inner">
          <Link to="/" className="logo" onClick={() => setMenuOpen(false)}>
            Hackathon<span>Hub</span>
          </Link>

          <button
            type="button"
            className={`menu-toggle ${menuOpen ? 'menu-toggle-open' : ''}`}
            aria-label={menuOpen ? 'Закрыть меню' : 'Открыть меню'}
            aria-expanded={menuOpen}
            onClick={() => setMenuOpen((v) => !v)}
          >
            <span />
            <span />
            <span />
          </button>

          <nav className={`nav glass ${menuOpen ? 'nav-open' : ''}`}>
            {user ? (
              <>
                <Link to="/profile" onClick={() => setMenuOpen(false)}>
                  Профиль
                </Link>
                {user.platform_role === 'organizer' && (
                  <>
                    <Link to="/organizer/hackathons" onClick={() => setMenuOpen(false)}>
                      Мои хакатоны
                    </Link>
                    <span className="nav-badge">организатор</span>
                  </>
                )}
                <button
                  type="button"
                  className="btn-ghost"
                  onClick={() => {
                    setMenuOpen(false);
                    logout();
                  }}
                >
                  Выйти
                </button>
              </>
            ) : (
              <>
                <Link to="/login" onClick={() => setMenuOpen(false)}>
                  Вход
                </Link>
                <Link to="/register" className="btn-primary btn-sm" onClick={() => setMenuOpen(false)}>
                  Регистрация
                </Link>
              </>
            )}
          </nav>
        </div>
      </header>

      <div
        className={`mobile-nav-backdrop ${menuOpen ? 'mobile-nav-backdrop-visible' : ''}`}
        onClick={() => setMenuOpen(false)}
        aria-hidden={!menuOpen}
      />

      <main className="main">
        <div key={location.pathname} className="page-enter">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
