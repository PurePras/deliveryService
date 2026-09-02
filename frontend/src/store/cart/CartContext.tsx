import { createContext, useContext, useEffect, useMemo, useReducer, type ReactNode } from 'react';
import type { Product } from '../../types';

const STORAGE_KEY = 'shri-ram-cart-v1';

export interface CartItem {
  productId: string;
  name: string;
  slug: string;
  price: string;
  unit: string;
  imageUrl?: string;
  isAvailable: boolean;
  quantity: number;
}

interface CartState {
  items: CartItem[];
}

type CartAction =
  | { type: 'HYDRATE'; items: CartItem[] }
  | { type: 'ADD_ITEM'; product: Product; quantity: number }
  | { type: 'SET_QUANTITY'; productId: string; quantity: number }
  | { type: 'REMOVE_ITEM'; productId: string }
  | { type: 'CLEAR_CART' };

function cartReducer(state: CartState, action: CartAction): CartState {
  switch (action.type) {
    case 'HYDRATE':
      return { items: action.items };

    case 'ADD_ITEM': {
      const existing = state.items.find((item) => item.productId === action.product.id);
      if (existing) {
        return {
          items: state.items.map((item) =>
            item.productId === action.product.id
              ? { ...item, quantity: item.quantity + action.quantity }
              : item,
          ),
        };
      }
      const { product } = action;
      const newItem: CartItem = {
        productId: product.id,
        name: product.name,
        slug: product.slug,
        price: product.price,
        unit: product.unit,
        imageUrl: product.image_url,
        isAvailable: product.is_available,
        quantity: action.quantity,
      };
      return { items: [...state.items, newItem] };
    }

    case 'SET_QUANTITY':
      if (action.quantity <= 0) {
        return { items: state.items.filter((item) => item.productId !== action.productId) };
      }
      return {
        items: state.items.map((item) =>
          item.productId === action.productId ? { ...item, quantity: action.quantity } : item,
        ),
      };

    case 'REMOVE_ITEM':
      return { items: state.items.filter((item) => item.productId !== action.productId) };

    case 'CLEAR_CART':
      return { items: [] };

    default:
      return state;
  }
}

function loadInitialState(): CartState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return { items: [] };
    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) return { items: [] };
    return { items: parsed as CartItem[] };
  } catch {
    return { items: [] };
  }
}

interface CartContextValue {
  items: CartItem[];
  itemCount: number;
  subtotal: number;
  addItem: (product: Product, quantity?: number) => void;
  setQuantity: (productId: string, quantity: number) => void;
  removeItem: (productId: string) => void;
  clearCart: () => void;
}

const CartContext = createContext<CartContextValue | undefined>(undefined);

export function CartProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(cartReducer, undefined, loadInitialState);

  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(state.items));
    } catch {
      // localStorage unavailable (private browsing, quota, etc.) — cart just won't persist.
    }
  }, [state.items]);

  const value = useMemo<CartContextValue>(() => {
    const itemCount = state.items.reduce((sum, item) => sum + item.quantity, 0);
    const subtotal = state.items.reduce((sum, item) => sum + Number(item.price) * item.quantity, 0);

    return {
      items: state.items,
      itemCount,
      subtotal,
      addItem: (product, quantity = 1) => dispatch({ type: 'ADD_ITEM', product, quantity }),
      setQuantity: (productId, quantity) => dispatch({ type: 'SET_QUANTITY', productId, quantity }),
      removeItem: (productId) => dispatch({ type: 'REMOVE_ITEM', productId }),
      clearCart: () => dispatch({ type: 'CLEAR_CART' }),
    };
  }, [state.items]);

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

export function useCart(): CartContextValue {
  const ctx = useContext(CartContext);
  if (!ctx) throw new Error('useCart must be used within a CartProvider');
  return ctx;
}
