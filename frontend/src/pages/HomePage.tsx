import { Link } from "react-router-dom";

import { useAuth } from "../context/AuthContext";

export function HomePage() {
  const { user } = useAuth();

  return (
    <section className="hero">
      <p className="eyebrow">Платформа хакатонов</p>
      <h1>Команды, кейсы и сдача проектов в одном месте</h1>
      <p className="lead">
        Регистрируйтесь на события, собирайте команду, выбирайте трек и сдавайте
        репозиторий, демо и питч до дедлайна.
      </p>

      {user ? (
        <div className="hero-actions">
          <p className="greeting">
            Привет, <strong>{user.full_name}</strong>
          </p>
          <Link to="/profile" className="btn-primary">
            Мой профиль
          </Link>
        </div>
      ) : (
        <div className="hero-actions">
          <Link to="/register" className="btn-primary">
            Создать аккаунт
          </Link>
          <Link to="/login" className="btn-secondary">
            Войти
          </Link>
        </div>
      )}

      <div className="card catalog-placeholder">
        <h2>Хакатоны</h2>
        <p>Каталог событий появится на следующем этапе разработки.</p>
      </div>
    </section>
  );
}
