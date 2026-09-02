import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../store/auth/AuthContext';
import { ApiError } from '../api/http';
import styles from './AuthForm.module.css';

function Register() {
  const { register } = useAuth();
  const navigate = useNavigate();

  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await register({ name, phone, password, email: email || undefined });
      navigate('/', { replace: true });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className={styles.page}>
      <form className={styles.form} onSubmit={handleSubmit}>
        <h1>Create Account</h1>
        <p className={styles.subtitle}>Join Shri Ram Sabji Delivery for fresh produce at your door.</p>

        {error && <p className={styles.error}>{error}</p>}

        <label className={styles.field}>
          <span>Full name</span>
          <input
            type="text"
            required
            autoComplete="name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="Your name"
          />
        </label>

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
          <span>Email (optional)</span>
          <input
            type="email"
            autoComplete="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            placeholder="you@example.com"
          />
        </label>

        <label className={styles.field}>
          <span>Password</span>
          <input
            type="password"
            required
            minLength={8}
            autoComplete="new-password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            placeholder="••••••••"
          />
          <span className={styles.hint}>At least 8 characters.</span>
        </label>

        <button type="submit" className={styles.submit} disabled={submitting}>
          {submitting ? 'Creating account…' : 'Create Account'}
        </button>

        <p className={styles.switch}>
          Already have an account? <Link to="/login">Log in</Link>
        </p>
      </form>
    </section>
  );
}

export default Register;
