const currencyFormatter = new Intl.NumberFormat('en-IN', {
  style: 'currency',
  currency: 'INR',
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
});

export function formatPrice(value: string | number): string {
  const amount = Number(value);
  return Number.isFinite(amount) ? currencyFormatter.format(amount) : String(value);
}

const unitLabels: Record<string, string> = {
  kg: 'kg',
  g: 'g',
  piece: 'piece',
  bundle: 'bundle',
  dozen: 'dozen',
};

export function formatPricePerUnit(price: string, unit: string): string {
  return `${formatPrice(price)} / ${unitLabels[unit] ?? unit}`;
}

const quantityFormatter = new Intl.NumberFormat('en-IN', {
  minimumFractionDigits: 0,
  maximumFractionDigits: 2,
});

export function formatQuantity(value: string): string {
  const amount = Number(value);
  return Number.isFinite(amount) ? quantityFormatter.format(amount) : value;
}
