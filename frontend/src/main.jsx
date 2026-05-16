import React from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Link, NavLink, Route, Routes } from "react-router-dom";
import { api, setAuthToken } from "./api";
import { entities, groups } from "./entities";
import { reports } from "./reports";
import "./styles.css";

const pages = [
  ["/", "Главная", "home"],
  ["/catalogs", "Справочники", "catalogs"],
  ["/employees", "Сотрудники", "employees"],
  ["/vehicles", "Транспорт", "vehicles"],
  ["/routes", "Маршруты", "routes"],
  ["/transport", "Перевозки", "transport"],
  ["/repairs", "Ремонт", "repairs"],
  ["/requests", "Комплектующие", "requests"],
  ["/reports", "Отчеты", "reports"],
];

const demoUsers = [
  ["director", "director_pass", "management"],
  ["workshop", "workshop_pass", "workshop_heads"],
  ["foreman", "foreman_pass", "foremen"],
  ["dispatcher", "dispatcher_pass", "dispatchers"],
  ["accountant", "accounting_pass", "accounting"],
  ["hr", "hr_pass", "hr"],
  ["driver", "driver_pass", "drivers_role"],
  ["repairman", "repairman_pass", "repairmen_role"],
];

function App() {
  const [session, setSession] = React.useState(null);
  const [loading, setLoading] = React.useState(Boolean(localStorage.getItem("token")));

  React.useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return;
    setAuthToken(token);
    api("/auth/me")
      .then((data) => setSession({ token, user: data.user, permissions: data.permissions }))
      .catch(() => {
        setAuthToken("");
        setSession(null);
      })
      .finally(() => setLoading(false));
  }, []);

  function login(next) {
    setAuthToken(next.token);
    setSession(next);
  }

  function logout() {
    setAuthToken("");
    setSession(null);
  }

  if (loading) return <main><p>Загрузка...</p></main>;
  if (!session) return <LoginPage onLogin={login} />;

  const permissions = session.permissions || session.user.permissions;
  const visiblePages = pages.filter(([, , key]) => hasPageAccess(key, permissions));

  return (
    <BrowserRouter>
      <header>
        <div className="topline">
          <h1>Информационная система автопредприятия</h1>
          <div className="userbar">
            {session.user.username} / {session.user.role}
            <button onClick={logout}>Выйти</button>
          </div>
        </div>
        <nav>
          {visiblePages.map(([to, title]) => (
            <NavLink key={to} to={to}>{title}</NavLink>
          ))}
        </nav>
      </header>
      <main>
        <Routes>
          <Route path="/" element={<Home permissions={permissions} />} />
          <Route path="/catalogs" element={<Group title="Справочники" names={groups.catalogs} permissions={permissions} />} />
          <Route path="/employees" element={<Group title="Сотрудники" names={groups.employees} permissions={permissions} />} />
          <Route path="/vehicles" element={<Group title="Транспорт" names={groups.vehicles} permissions={permissions} />} />
          <Route path="/routes" element={<Group title="Маршруты и назначения" names={groups.routes} permissions={permissions} />} />
          <Route path="/transport" element={<Group title="Перевозки" names={groups.transport} permissions={permissions} />} />
          <Route path="/repairs" element={<Group title="Ремонт" names={groups.repairs} permissions={permissions} />} />
          <Route path="/requests" element={<Group title="Комплектующие" names={groups.requests} permissions={permissions} />} />
          <Route path="/reports" element={<ReportsPage permissions={permissions} />} />
        </Routes>
      </main>
    </BrowserRouter>
  );
}

function LoginPage({ onLogin }) {
  const [form, setForm] = React.useState({ username: "director", password: "director_pass" });
  const [error, setError] = React.useState("");

  async function submit(e) {
    e.preventDefault();
    try {
      const data = await api("/auth/login", { method: "POST", body: JSON.stringify(form) });
      setError("");
      onLogin(data);
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <main className="login-page">
      <section className="login-box">
        <h1>Вход</h1>
        {error && <div className="error">{error}</div>}
        <form onSubmit={submit} className="entity-form">
          <label>
            <span>Логин</span>
            <input value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} />
          </label>
          <label>
            <span>Пароль</span>
            <input type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} />
          </label>
          <button>Войти</button>
        </form>
        <table>
          <thead><tr><th>Логин</th><th>Пароль</th><th>Роль</th></tr></thead>
          <tbody>
            {demoUsers.map(([username, password, role]) => (
              <tr key={username} onClick={() => setForm({ username, password })}>
                <td>{username}</td><td>{password}</td><td>{role}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    </main>
  );
}

function Home({ permissions }) {
  return (
    <section>
      <h2>Главная</h2>
      <p>Backend: Go + chi + pgx. Database: PostgreSQL. Frontend: React.</p>
      {Object.keys(permissions.reports || {}).length > 0 && <p><Link to="/reports">Перейти к отчетам</Link></p>}
    </section>
  );
}

function Group({ title, names, permissions }) {
  const visible = names.filter((name) => canEntity(permissions, name, "read"));
  return (
    <section>
      <h2>{title}</h2>
      {visible.length === 0 && <p className="muted">Нет доступных разделов для текущей роли</p>}
      {visible.map((name) => <EntityTable key={name} cfg={entities[name]} permissions={permissions} />)}
    </section>
  );
}

function EntityTable({ cfg, permissions }) {
  const idCol = cfg.id || "id";
  const canCreate = canEntity(permissions, cfg.endpoint, "create");
  const canUpdate = canEntity(permissions, cfg.endpoint, "update");
  const canDelete = canEntity(permissions, cfg.endpoint, "delete");
  const [rows, setRows] = React.useState([]);
  const [form, setForm] = React.useState(emptyForm(cfg));
  const [editId, setEditId] = React.useState(null);
  const [error, setError] = React.useState("");

  const load = React.useCallback(async () => {
    try {
      setRows(await api(`/${cfg.endpoint}`));
      setError("");
    } catch (e) {
      setError(e.message);
    }
  }, [cfg.endpoint]);

  React.useEffect(() => { load(); }, [load]);

  async function submit(e) {
    e.preventDefault();
    try {
      const body = JSON.stringify(prepare(cfg, form));
      if (editId == null) {
        await api(`/${cfg.endpoint}`, { method: "POST", body });
      } else {
        await api(`/${cfg.endpoint}/${editId}`, { method: "PUT", body });
      }
      setForm(emptyForm(cfg));
      setEditId(null);
      await load();
    } catch (e) {
      setError(e.message);
    }
  }

  async function remove(id) {
    if (!window.confirm("Удалить запись?")) return;
    try {
      await api(`/${cfg.endpoint}/${id}`, { method: "DELETE" });
      await load();
    } catch (e) {
      setError(e.message);
    }
  }

  function edit(row) {
    const next = {};
    cfg.fields.forEach(([name]) => { next[name] = row[name] ?? ""; });
    setForm(next);
    setEditId(row[idCol]);
  }

  return (
    <section className="entity">
      <h3>{cfg.title}</h3>
      {error && <div className="error">{error}</div>}
      {(canCreate || (editId != null && canUpdate)) && (
        <form onSubmit={submit} className="entity-form">
          {cfg.fields.map(([name, label, type = "text", nullable]) => (
            <label key={name}>
              <span>{label}</span>
              <input
                type={type === "number" ? "number" : type === "date" ? "date" : "text"}
                step={type === "number" ? "0.01" : undefined}
                value={form[name] ?? ""}
                required={!nullable}
                disabled={editId != null && name === idCol}
                onChange={(e) => setForm({ ...form, [name]: e.target.value })}
              />
            </label>
          ))}
          <button>{editId == null ? "Создать" : "Сохранить"}</button>
          {editId != null && <button type="button" onClick={() => { setEditId(null); setForm(emptyForm(cfg)); }}>Отмена</button>}
        </form>
      )}
      <DataTable
        rows={rows}
        columns={[idCol, ...cfg.fields.map(([name]) => name).filter((name) => name !== idCol)]}
        onEdit={canUpdate ? edit : null}
        onDelete={canDelete ? (row) => remove(row[idCol]) : null}
      />
    </section>
  );
}

function emptyForm(cfg) {
  const data = {};
  cfg.fields.forEach(([name]) => { data[name] = ""; });
  return data;
}

function prepare(cfg, form) {
  const body = {};
  cfg.fields.forEach(([name, , type = "text", nullable]) => {
    const value = form[name];
    if (value === "" && nullable) {
      body[name] = null;
    } else if (value === "" && type !== "text") {
      body[name] = null;
    } else if (type === "number") {
      body[name] = Number(value);
    } else {
      body[name] = value;
    }
  });
  return body;
}

function DataTable({ rows, columns, onEdit, onDelete }) {
  const cols = columns || Array.from(new Set(rows.flatMap((r) => Object.keys(r))));
  if (!rows.length) return <p className="muted">Нет данных</p>;
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>
            {cols.map((c) => <th key={c}>{c}</th>)}
            {(onEdit || onDelete) && <th>Действия</th>}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, idx) => (
            <tr key={idx}>
              {cols.map((c) => <td key={c}>{format(row[c])}</td>)}
              {(onEdit || onDelete) && (
                <td className="actions">
                  {onEdit && <button onClick={() => onEdit(row)}>Изм.</button>}
                  {onDelete && <button onClick={() => onDelete(row)}>Удалить</button>}
                </td>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function ReportsPage({ permissions }) {
  const visible = reports.filter(([slug]) => permissions.reports?.[slug]);
  return (
    <section>
      <h2>Отчеты</h2>
      {visible.length === 0 && <p className="muted">Нет доступных отчетов для текущей роли</p>}
      {visible.map(([slug, title, params]) => (
        <ReportBlock key={slug} slug={slug} title={title} params={params} />
      ))}
    </section>
  );
}

function ReportBlock({ slug, title, params }) {
  const [values, setValues] = React.useState(defaultParams(params));
  const [rows, setRows] = React.useState([]);
  const [error, setError] = React.useState("");

  async function run(e) {
    e.preventDefault();
    const qs = new URLSearchParams();
    Object.entries(values).forEach(([key, value]) => {
      if (value !== "") qs.set(key, value);
    });
    try {
      setRows(await api(`/reports/${slug}${qs.toString() ? `?${qs}` : ""}`));
      setError("");
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <section className="report">
      <h3>{title}</h3>
      <form onSubmit={run} className="report-form">
        {params.map(([name, label, type = "text"]) => (
          <label key={name}>
            <span>{label}</span>
            <input
              type={type}
              value={values[name] || ""}
              onChange={(e) => setValues({ ...values, [name]: e.target.value })}
            />
          </label>
        ))}
        <button>Получить</button>
      </form>
      {error && <div className="error">{error}</div>}
      <DataTable rows={rows} />
    </section>
  );
}

function defaultParams(params) {
  const values = {};
  params.forEach(([name]) => { values[name] = ""; });
  return values;
}

function hasPageAccess(key, permissions) {
  if (key === "home") return true;
  if (key === "reports") return Object.values(permissions.reports || {}).some(Boolean);
  return (groups[key] || []).some((name) => canEntity(permissions, name, "read"));
}

function canEntity(permissions, entity, action) {
  return (permissions.entities?.[entity] || []).includes(action);
}

function format(value) {
  if (value == null) return "";
  if (typeof value === "number") return Number.isInteger(value) ? value : value.toFixed(2);
  return String(value);
}

createRoot(document.getElementById("root")).render(<App />);
