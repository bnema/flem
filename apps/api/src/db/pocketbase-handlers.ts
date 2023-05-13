import PocketBase from 'pocketbase';
import dotenv from 'dotenv';

dotenv.config();

class PocketBaseService {
  pb: PocketBase;
  constructor() {
    this.pb = new PocketBase('https://pocketbase.flem.bnema.dev');
  }

  async authenticate() {
    try {
      const authData = await this.pb.admins.authWithPassword(
        process.env.POCKETBASE_EMAIL || '',
        process.env.POCKETBASE_PASSWORD || ''
      );
      return authData;
    } catch (error) {
      console.error('Error during authentication', error);
      return null;
    }
  }

  logout() {
    this.pb.authStore.clear();
  }

  async getAllUsers() {
    const authData = await this.authenticate();
    if (authData) {
      try {
        const users = await this.pb.collection('users').getFullList();
        return users;
      } catch (error) {
        console.error('Error during getting users', error);
        return null;
      }
    } else {
      console.error('Error during authentication');
      return null;
    }
  }
}

export const pocketBaseService = new PocketBaseService();
