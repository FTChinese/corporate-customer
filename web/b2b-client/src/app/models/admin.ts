// Credentials is the login or signup
// data.
export interface Credentials {
  email: string;
  password: string;
}

export interface Account {
  id: string;
  email: string;
  displayName: string | null;
  active: boolean;
  verified: boolean;
}

// Passport is what returend from API after user logged in
// or signed up.
export interface Passport extends Account {
  teamId: string | null; // null indicates team is not created.
  expiresAt: number;
  token: string; // Json web token.
}

export interface Profile extends Account {
  createdUtc: string | null;
  updatedUtc: string | null;
}

export interface Passwords {
  oldPassword: string;
  password: string;
}
