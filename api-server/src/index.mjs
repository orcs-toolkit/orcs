import dotenv from 'dotenv';
dotenv.config();
import express from 'express';
import cors from 'cors';
import passport from 'passport';
import helmet from 'helmet';
import mongoose from 'mongoose';

import { morganMiddleware } from './middlewares/morganMiddleware.mjs';
import { sessionConfig } from './services/sessionStore.mjs';
import { routes } from './routes/index.mjs';
import { Logger } from './services/logger.mjs';
const logger = new Logger();

//passport config file
import('./services/passportConfig.mjs');

const app = express();
app.use(express.json());
app.use(
	cors({
		origin: true,
		credentials: true,
	})
);
app.use(sessionConfig);
app.use(passport.initialize());
app.use(passport.session());
app.use(helmet());
app.use(morganMiddleware);

// Routes
app.use(routes);

(async () => {
	try {
		await mongoose.connect(process.env.MONGO_URI);
		logger.info('Connected to MongoDB');
	} catch (err) {
		logger.error(err);
	}
})();

const PORT = 4001;
app.listen(process.env.PORT || PORT, () => {
	logger.info(`Auth server running on port: ${PORT}`);
});