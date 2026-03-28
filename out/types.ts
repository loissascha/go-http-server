export interface loginInput {
  username: string;
  password: string;
}
export interface loginResult {
  method: string;
  success: boolean;
  jwt: string;
}
