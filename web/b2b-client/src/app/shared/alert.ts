export type AlertType = 'success' | 'info' | 'warning' | 'danger' | 'primary' | 'secondary' | 'light' | 'dark';

export interface Alert {
  type: AlertType;
  message: string;
}
