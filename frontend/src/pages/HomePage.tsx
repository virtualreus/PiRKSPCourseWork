import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import * as hackathonsApi from "../api/hackathons";
import type { HackathonListItem } from "../api/hackathonTypes";
import { ApiError } from "../api/client";
import { HackathonCard } from "../components/HackathonCard";
import { Reveal } from "../components/Reveal";
import { useAuth } from "../context/AuthContext";

const flowSteps = [
  {
    step: "01",
    title: "Запишитесь на ивент",
    text: "Создайте аккаунт участника и подтвердите участие в хакатоне, пока открыта регистрация.",
  },
  {
    step: "02",
    title: "Соберите команду и кейс",
    text: "Создайте команду или вступите в существующую. Выберите трек и задачу от заказчика.",
  },
  {
    step: "03",
    title: "Кодьте до дедлайна",
    text: "Работайте в своём темпе в окне хакатона. Все даты и статусы - на платформе.",
  },
  {
    step: "04",
    title: "Сдайте артефакты",
    text: "Репозиторий, демо, питч - одним сабмитом до submission_deadline.",
  },
];

const features = [
  {
    title: "Треки и кейсы",
    text: "Реальные задачи от заказчиков с ресурсами и описанием - как на ЛЦТ и Цифровом прорыве.",
    accent: "cyan",
  },
  {
    title: "Прозрачный таймлайн",
    text: "Регистрация, старт кодинга и дедлайн сдачи - без сюрпризов для команд и жюри.",
    accent: "blue",
  },
  {
    title: "Роли без хаоса",
    text: "Организатор ведёт событие, участник проходит путь от записи до сабмита.",
    accent: "purple",
  },
];

export function HomePage() {
  const { user } = useAuth();
  const [items, setItems] = useState<HackathonListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    (async () => {
      try {
        const resp = await hackathonsApi.listHackathons();
        setItems(resp.items ?? []);
      } catch (err) {
        if (err instanceof ApiError) {
          setError(err.message);
        } else {
          setError("Не удалось загрузить каталог");
        }
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const openCount = items.filter((h) => h.status === "registration").length;

  return (
    <div className="landing">
      <section className="landing-hero">
        <div className="landing-hero-glow" aria-hidden />
        <div className="stagger">
          <Reveal>
            <p className="landing-eyebrow">
              <span className="pulse-dot" />
              HackathonHub · платформа соревнований
            </p>
          </Reveal>
          <Reveal delay={80}>
            <h1 className="landing-title">
              Собери команду.
              <br />
              <span className="text-gradient">Взорви кейс.</span>
              <br />
              Сдай до дедлайна.
            </h1>
          </Reveal>
          <Reveal delay={160}>
            <p className="landing-lead">
              Единая среда для хакатонов уровня{" "}
              <strong>«Лидеры цифровой трансформации»</strong> и{" "}
              <strong>«Цифровой прорыв»</strong>: каталог событий, треки, кейсы,
              команды и артефакты сдачи - без Excel и хаоса в чатах.
            </p>
          </Reveal>
          <Reveal delay={220}>
            <div className="landing-hero-actions">
              {user ? (
                <>
                  <p className="landing-greeting">
                    С возвращением, <strong>{user.full_name}</strong>
                  </p>
                  <a href="#catalog" className="btn-primary">
                    Смотреть хакатоны
                  </a>
                  {user.platform_role === "organizer" ? (
                    <Link to="/organizer/hackathons" className="btn-secondary">
                      Кабинет организатора
                    </Link>
                  ) : (
                    <Link to="/profile" className="btn-secondary">
                      Профиль
                    </Link>
                  )}
                </>
              ) : (
                <>
                  <Link to="/register" className="btn-primary">
                    Начать бесплатно
                  </Link>
                  <Link to="/login" className="btn-secondary">
                    Уже есть аккаунт
                  </Link>
                  <a href="#catalog" className="btn-ghost">
                    Каталог ↓
                  </a>
                </>
              )}
            </div>
          </Reveal>
        </div>

        <Reveal className="landing-chips" delay={280}>
          <span className="landing-chip glass">
            {items.length || "-"} событий в каталоге
          </span>
          <span className="landing-chip glass">
            {openCount || "-"} с открытой регистрацией
          </span>
          <span className="landing-chip glass">48–72 ч кодинг</span>
          <span className="landing-chip glass">repo · demo · pitch</span>
        </Reveal>
      </section>

      <div className="landing-marquee" aria-hidden>
        <div className="landing-marquee-track">
          <span>кейсы</span>
          <span>команды</span>
          <span>треки</span>
          <span>демо</span>
          <span>питч</span>
          <span>дедлайн</span>
          <span>open data</span>
          <span>прототип</span>
          <span>кейсы</span>
          <span>команды</span>
          <span>треки</span>
          <span>демо</span>
          <span>питч</span>
          <span>дедлайн</span>
        </div>
      </div>

      <Reveal>
        <section className="landing-flow">
          <div className="landing-section-head">
            <p className="eyebrow">Как это работает</p>
            <h2>От идеи до сдачи - четыре шага</h2>
          </div>
          <div className="landing-flow-grid">
            {flowSteps.map((item, i) => (
              <Reveal key={item.step} delay={i * 70}>
                <article className="landing-flow-card glass">
                  <span className="landing-flow-step">{item.step}</span>
                  <h3>{item.title}</h3>
                  <p>{item.text}</p>
                </article>
              </Reveal>
            ))}
          </div>
        </section>
      </Reveal>

      <Reveal>
        <section className="landing-features">
          <div className="landing-section-head">
            <p className="eyebrow">Почему HackathonHub</p>
            <h2>Сделано под настоящие хакатоны</h2>
          </div>
          <div className="landing-features-grid">
            {features.map((f, i) => (
              <Reveal key={f.title} delay={i * 80}>
                <article className={`landing-feature glass accent-${f.accent}`}>
                  <h3>{f.title}</h3>
                  <p>{f.text}</p>
                </article>
              </Reveal>
            ))}
          </div>
        </section>
      </Reveal>

      <Reveal>
        <section className="landing-banner glass">
          <div className="landing-banner-inner">
            <div>
              <p className="eyebrow">Для организаторов</p>
              <h2>Публикуйте хакатон за 10 минут</h2>
              <p className="landing-banner-text">
                Черновик → треки и кейсы → публикация. Участники видят карточку
                в каталоге сразу после publish.
              </p>
            </div>
            <div className="landing-banner-actions">
              {user?.platform_role === "organizer" ? (
                <Link to="/organizer/hackathons/new" className="btn-primary">
                  Создать хакатон
                </Link>
              ) : (
                <Link to="/register" className="btn-primary">
                  Стать организатором
                </Link>
              )}
              <p className="landing-demo-hint">
                Демо: organizer@demo.local / demo12345
              </p>
            </div>
          </div>
        </section>
      </Reveal>

      <Reveal>
        <section id="catalog" className="catalog-section landing-catalog">
          <div className="landing-section-head landing-section-head-row">
            <div>
              <p className="eyebrow">Каталог</p>
              <h2>Актуальные хакатоны</h2>
            </div>
            {!loading && items.length > 0 && (
              <span className="landing-catalog-count glass">
                {items.length} событий
              </span>
            )}
          </div>

          {loading && (
            <div className="skeleton-grid" aria-busy>
              <div className="skeleton-card glass" />
              <div className="skeleton-card glass" />
              <div className="skeleton-card glass" />
            </div>
          )}
          {error && <p className="form-error">{error}</p>}
          {!loading && (
            <div className="hackathon-grid">
              {items.map((item, i) => (
                <Reveal key={item.id} delay={i * 60}>
                  <HackathonCard item={item} />
                </Reveal>
              ))}
            </div>
          )}
        </section>
      </Reveal>
    </div>
  );
}
