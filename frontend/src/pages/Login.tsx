import { useEffect, useState, type FormEvent } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../store/auth/AuthContext';
import { ApiError } from '../api/http';
import styles from './AuthForm.module.css';

const RESEND_COOLDOWN_SECONDS = 30;
const PHONE_PATTERN = /^[6-9]\d{9}$/;
const OTP_PATTERN = /^\d{6}$/;

// Matches the example in the spec: all but the last 4 digits replaced with '*'.
function maskPhone(phone: string): string {
  return phone.length <= 4 ? phone : '*'.repeat(phone.length - 4) + phone.slice(-4);
}

function Login() {
  const { requestOtp, verifyOtp } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [step, setStep] = useState<'phone' | 'otp'>('phone');
  const [phone, setPhone] = useState('');
  const [code, setCode] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [cooldown, setCooldown] = useState(0);

  const from = (location.state as { from?: string } | null)?.from ?? '/';

  useEffect(() => {
    if (cooldown <= 0) return;
    const timer = window.setInterval(() => setCooldown((seconds) => Math.max(0, seconds - 1)), 1000);
    return () => window.clearInterval(timer);
  }, [cooldown]);

  async function handleSendOtp(event: FormEvent) {
    event.preventDefault();
    setError(null);
    if (!PHONE_PATTERN.test(phone)) {
      setError('Enter a valid 10-digit mobile number.');
      return;
    }
    setSubmitting(true);
    try {
      await requestOtp(phone);
      setStep('otp');
      setCode('');
      setCooldown(RESEND_COOLDOWN_SECONDS);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.');
    } finally {
      setSubmitting(false);
    }
  }

  async function handleVerify(event: FormEvent) {
    event.preventDefault();
    setError(null);
    if (!OTP_PATTERN.test(code)) {
      setError('Enter the 6-digit code.');
      return;
    }
    setSubmitting(true);
    try {
      await verifyOtp(phone, code);
      navigate(from, { replace: true });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.');
    } finally {
      setSubmitting(false);
    }
  }

  async function handleResend() {
    setError(null);
    setSubmitting(true);
    try {
      await requestOtp(phone);
      setCode('');
      setCooldown(RESEND_COOLDOWN_SECONDS);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Please try again.');
    } finally {
      setSubmitting(false);
    }
  }

  function handleChangeNumber() {
    setStep('phone');
    setCode('');
    setError(null);
    setCooldown(0);
  }

  if (step === 'otp') {
    return (
      <section className={styles.page}>
        <form className={styles.form} onSubmit={handleVerify}>
          <h1>Enter OTP</h1>
          <p className={styles.subtitle}>OTP sent to {maskPhone(phone)}</p>

          {error && <p className={styles.error}>{error}</p>}

          <label className={styles.field}>
            <span>6-digit code</span>
            <input
              type="text"
              inputMode="numeric"
              pattern="\d{6}"
              maxLength={6}
              required
              autoComplete="one-time-code"
              autoFocus
              value={code}
              onChange={(event) => setCode(event.target.value.replace(/\D/g, '').slice(0, 6))}
              placeholder="123456"
            />
          </label>

          <button type="submit" className={styles.submit} disabled={submitting}>
            {submitting ? 'Verifying…' : 'Verify OTP'}
          </button>

          <div className={styles.otpActions}>
            {cooldown > 0 ? (
              <span>Resend OTP in {cooldown}s</span>
            ) : (
              <button type="button" className={styles.linkButton} onClick={handleResend} disabled={submitting}>
                Resend OTP
              </button>
            )}
            <span>·</span>
            <button type="button" className={styles.linkButton} onClick={handleChangeNumber} disabled={submitting}>
              Change mobile number
            </button>
          </div>
        </form>
      </section>
    );
  }

  return (
    <section className={styles.page}>
      <form className={styles.form} onSubmit={handleSendOtp}>
        <h1>Log In</h1>
        <p className={styles.subtitle}>Welcome back to Shri Ram Sabji Delivery.</p>

        {error && <p className={styles.error}>{error}</p>}

        <label className={styles.field}>
          <span>Phone number</span>
          <input
            type="tel"
            inputMode="numeric"
            required
            autoComplete="tel"
            autoFocus
            value={phone}
            onChange={(event) => setPhone(event.target.value.replace(/\D/g, '').slice(0, 10))}
            placeholder="98765 43210"
          />
        </label>

        <button type="submit" className={styles.submit} disabled={submitting}>
          {submitting ? 'Sending…' : 'Send OTP'}
        </button>

        <p className={styles.switch}>
          New here? <Link to="/register">Create an account</Link>
        </p>
      </form>
    </section>
  );
}

export default Login;
