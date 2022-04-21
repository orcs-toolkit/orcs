import { useEffect, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';

import axios from 'axios';

export default function SessionDashboard() {
	const [user, setUser] = useState({});
	const [policy, setPolicy] = useState([]);
	const [showList, setShowList] = useState(false);
	const [sessionStarted, setSession] = useState('');
	const [error, setError] = useState('');
	const history = useHistory();

	useEffect(() => {
		const token = localStorage.getItem('jwt');

		(async () => {
			const { data } = await axios.get(
				'http://localhost:4001/auth/currentuser',
				{
					headers: {
						'Content-Type': 'application/json',
						Authorization: `Bearer ${token}`,
					},
				}
			);
			setUser(data.currentUser);
		})();
	}, []);

	async function handleSession(e) {
		e.preventDefault();

		const payload = {
			name: user.name,
			role: user.role,
		};

		try {
			const { data } = await axios.post('http://localhost:4002/role', payload, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Allow-Origin-With-Credentials': '*',
					'Access-Control-Allow-Origin': '*',
				},
			});
			// console.log(typeof data.success);
			if (!data.success) {
				setError('Something went wrong, please try again');
			}

			setSession('Session started!');
			setError('');
		} catch (err) {
			console.log('Unable to POST payload', err);
			setError('Something went wrong, please try again');
			setSession('');
		}
	}

	async function handleLogout(e) {
		e.preventDefault();

		localStorage.removeItem('jwt');

		const payload = {
			name: '',
			role: 'default',
		};

		try {
			const { data } = await axios.post(
				'http://localhost:4002/logout',
				payload,
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						'Allow-Origin-With-Credentials': '*',
						'Access-Control-Allow-Origin': '*',
					},
				}
			);
			// console.log(typeof data.success);
			if (!data.success) {
				setSession('Something went wrong, please try again');
			}

			setSession('');
		} catch (err) {
			console.log('Unable to POST payload', err);
		}

		history.push('/');
	}

	async function fetchPolicy(e) {
		e.preventDefault();

		const { data } = await axios.get(
			`http://localhost:4001/policy/getRolePolicy/${user.role}`
		);
		if (!data.policy) {
			return;
		}
		setPolicy(data.policy);
	}

	function getPolicyList() {
		console.log('Policy', policy);
		if (policy.length === 0) {
			return;
		}
		if (policy.banList.length === 0) {
			return <p className="list-group-item h5">No policies set</p>;
		}
		return policy.banList.map((i) => {
			return <li className="list-group-item">{i}</li>;
		});
	}

	return (
		<section className="vh-100" style={{ backgroundColor: '#eee' }}>
			<div className="container h-100">
				<div className="row justify-content-center align-items-center h-100">
					<div className="col-lg-12 col-xl-11">
						<div className="card text-black" style={{ borderRadius: '25px' }}>
							<div className="card-body p-md-5">
								<div className="row justify-content-center">
									<div className="col-md-10 col-lg-6 col-xl-5 order-2 order-lg-1">
										<p className="text-center h1 fw-bold mb-5 mx-1 mx-md-4 mt-4">
											Welcome,{' '}
											<span className="text-capitalize">{user.name}</span>
										</p>

										<div>
											<h3>Please check your information</h3>
											<p className="fst-normal fs-4">
												Name:{' '}
												<span className="text-capitalize">{user.name}</span>
												<br />
												Email: {user.email}
												<br />
												Role: {user.role}
											</p>
										</div>
										<div style={{ marginTop: '10px' }} className="mb-4">
											Something wrong? &nbsp;
											<Link to="/login">Login again!</Link>
										</div>

										<div className="row">
											<div className="col-sm-6">
												<form onSubmit={handleSession} className="mx-1 mx-md-4">
													<div className="d-flex mb-lg-4">
														<button
															type="submit"
															className="btn btn-primary btn-lg"
														>
															Start your session!
														</button>
													</div>
												</form>
											</div>

											<div className="col-sm-6">
												<form onSubmit={handleLogout} className="mx-1 mx-md-4">
													<div className="d-flex mb-lg-4">
														<button
															type="submit"
															className="btn btn-danger btn-lg"
														>
															Logout
														</button>
													</div>
												</form>
											</div>
										</div>

										<div className="row">
											<div className="col">
												<h4>Click here to know what you cannot use</h4>

												<button
													onClick={(e) => {
														fetchPolicy(e);
														setShowList(true);
													}}
													type="submit"
													className="btn btn-secondary btn-lg"
												>
													Fetch Policy
												</button>

												{showList && (
													<div className="card-text">
														<ul
															className="list-group list-group-flush list-group-numbered"
															style={{
																maxHeight: '200px',
																marginBottom: '10px',
																overflowY: 'scroll',
																WebkitOverflowScrolling: 'touch',
															}}
														>
															{getPolicyList()}
														</ul>
													</div>
												)}
											</div>
										</div>

										<br />

										{sessionStarted && (
											<div className="alert alert-success" role="alert">
												{sessionStarted}
											</div>
										)}
										{error && (
											<div className="alert alert-danger" role="alert">
												{error}
											</div>
										)}
									</div>
									<div className="col-md-10 col-lg-6 col-xl-7 d-flex align-items-center order-1 order-lg-2">
										<input
											type="image"
											src="https://mdbcdn.b-cdn.net/img/Photos/new-templates/bootstrap-registration/draw1.webp"
											className="img-fluid"
											alt="Sample image"
										/>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</section>
	);
}
