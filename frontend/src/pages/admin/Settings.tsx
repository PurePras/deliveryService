import { useState, type FormEvent } from 'react';
import { useAuth } from '../../store/auth/AuthContext';
import { changePassword } from '../../api/auth';
import { ApiError } from '../../api/http';
import styles from './Settings.module.css';

function AdminSettings() {
  const { user } = useAuth();
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  if (!user) return null;

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    setError(null);
    setSuccess(false);

    if (newPassword !== confirmPassword) {
      setError('New password and confirmation do not match.');
      return;
    }

    setSubmitting(true);
    try {
      await changePassword({ current_password: currentPassword, new_password: newPassword });
      setSuccess(true);
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className={styles.page}>
      <h1>Settings</h1>

      <div className={styles.card}>
        <h2>Profile</h2>
        <div className={styles.row}>
          <span className={styles.label}>Name</span>
          <span>{user.name}</span>
        </div>
        <div className={styles.row}>
          <span className={styles.label}>Phone</span>
          <span>{user.phone}</span>
        </div>
        {user.email && (
          <div className={styles.row}>
            <span className={styles.label}>Email</span>
            <span>{user.email}</span>
          </div>
        )}
        <div className={styles.row}>
          <span className={styles.label}>Role</span>
          <span className={styles.roleBadge}>{user.role}</span>
        </div>
      </div>

      <div className={styles.card}>
        <h2>Change Password</h2>
        <form className={styles.form} onSubmit={handleSubmit}>
          {error && <p className={styles.error}>{error}</p>}
          {success && <p className={styles.success}>Password changed successfully.</p>}

          <label className={styles.field}>
            <span>Current password</span>
            <input
              type="password"
              required
              autoComplete="current-password"
              value={currentPassword}
              onChange={(event) => setCurrentPassword(event.target.value)}
            />
          </label>

          <label className={styles.field}>
            <span>New password</span>
            <input
              type="password"
              required
              minLength={8}
              autoComplete="new-password"
              value={newPassword}
              onChange={(event) => setNewPassword(event.target.value)}
            />
          </label>

          <label className={styles.field}>
            <span>Confirm new password</span>
            <input
              type="password"
              required
              autoComplete="new-password"
              value={confirmPassword}
              onChange={(event) => setConfirmPassword(event.target.value)}
            />
          </label>

          <button type="submit" className={styles.submitBtn} disabled={submitting}>
            {submitting ? 'Saving…' : 'Change Password'}
          </button>
        </form>
      </div>
    </section>
  );
}

export default AdminSettings;
