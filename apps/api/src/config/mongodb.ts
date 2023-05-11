import dotenv from "dotenv";
// Load token from .env.local
dotenv.config();

// MongoDB config
import { connect } from 'mongoose';

export const connectDB = async () => {
    try {
        await connect(process.env.MONGODB_URL|| '');
    } catch (error) {
        console.error('Error in DB Connection: ' + error);
        process.exit(1);
    }
};