import { useState, type FormEvent } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../store/auth/AuthContext';
import { ApiError } from '../api/http';
import styles from './AuthForm.module.css';

function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const from = (location.state as { from?: string } | null)?.from ?? '/';

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await login({ phone, password });
      navigate(from, { replace: true });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className={styles.page}>
      <form className={styles.form} onSubmit={handleSubmit}>
        <h1>Log In</h1>
        <p className={styles.subtitle}>Welcome back to Shri Ram Sabji Delivery.</p>

        {error && <p className={styles.error}>{error}</p>}

        <label className={styles.field}>
          <span>Phone number</span>
          <input
            type="tel"
            required
            autoComplete="tel"
            value={phone}
            onChange={(event) => setPhone(event.target.value)}
            placeholder="98765 43210"
          />
        </label>

        <label className={styles.field}>
          <span>Password</span>
          <input
            type="password"
            required
            autoComplete="current-password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            placeholder="••••••••"
          />
        </label>

        <button type="submit" className={styles.submit} disabled={submitting}>
          {submitting ? 'Logging in…' : 'Log In'}
        </button>

        <p className={styles.switch}>
          New here? <Link to="/register">Create an account</Link>
        </p>
      </form>
    </section>
  );
}

export default Login;
