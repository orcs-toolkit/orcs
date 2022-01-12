const mongoose = require('mongoose');
const { Machine } = require('../models/Machine');

(async () => {
	if (!process.env.SECRET_KEY) {
		throw new Error('SECRET_KEY not defined!');
	}

	if (!process.env.MONGO_URI) {
		throw new Error('MONGO_URI not specified!');
	}

	try {
		await mongoose.connect(process.env.MONGO_URI);
		console.log('connected to MongoDB');
	} catch (err) {
		throw new Error(err);
	}
})();

function socketMain(io, socket) {
	let macA;

	console.log('Socket Session: ', socket.request.session);

	socket.on('clientAuth', (key) => {
		if (key === 'clientAuth') {
			// valid UI client joined
			socket.join('ui');
			console.log('React client joined!');
			Machine.find({}, (err, docs) => {
				docs.forEach((aMachine) => {
					// on load assume all machines are offline
					aMachine.isActive = false;
					io.to('ui').emit('data', aMachine);
				});
			});
		} else {
			// invalid client
			socket.disconnect(true);
		}
	});

	socket.on('disconnect', () => {
		Machine.find({ macA }, (err, docs) => {
			if (docs.length > 0) {
				// send one last emit to client
				docs[0].isActive = false;
				io.to('ui').emit('data', docs[0]);
			}
		});
	});

	socket.on('initPerfData', async (data) => {
		macA = data.macA;
		const mongooseResponse = await checkAndAdd(data);
		console.log(mongooseResponse);
	});

	socket.on('perfData', (data) => {
		io.to('ui').emit('data', data);
	});
}

function checkAndAdd(data) {
	return new Promise((resolve, reject) => {
		Machine.findOne({ macA: data.macA }, (err, docs) => {
			if (err) {
				throw err;
				// reject(err);
			} else if (docs === null) {
				// record or machine not added to db, so add it!
				let newMachine = new Machine(data);
				newMachine.save();
				resolve('Machine registered!');
			} else {
				// record or machine already added to db, so resolve
				resolve('Already registered!');
			}
		});
	});
}

module.exports = socketMain;
