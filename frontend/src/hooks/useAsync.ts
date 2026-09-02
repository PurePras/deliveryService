import { useEffect, useState } from 'react';

interface AsyncState<T> {
  data: T | undefined;
  loading: boolean;
  error: Error | undefined;
}

/**
 * Runs `fn` whenever `deps` change, tracking loading/error state and ignoring
 * results from a run that's been superseded by a newer one (stale closures
 * from fast filter/search changes).
 */
export function useAsync<T>(fn: () => Promise<T>, deps: unknown[]): AsyncState<T> {
  const [state, setState] = useState<AsyncState<T>>({ data: undefined, loading: true, error: undefined });

  useEffect(() => {
    let cancelled = false;
    setState((prev) => ({ ...prev, loading: true, error: undefined }));

    fn()
      .then((data) => {
        if (!cancelled) setState({ data, loading: false, error: undefined });
      })
      .catch((error: unknown) => {
        if (!cancelled) setState({ data: undefined, loading: false, error: error as Error });
      });

    return () => {
      cancelled = true;
    };
  }, deps);

  return state;
}
