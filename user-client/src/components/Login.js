import { useState } from 'react';
import { Link } from 'react-router-dom';

export default function Login() {
	const [dropDownValue, setDropDownValue] = useState('');

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
											Login
										</p>

										<form className="mx-1 mx-md-4">
											<div className="d-flex flex-row align-items-center mb-4">
												<i className="fas fa-user fa-lg me-3 fa-fw"></i>
												<div className="form-outline flex-fill mb-0">
													<input
														type="text"
														id="form3Example1c"
														className="form-control"
													/>
													<label
														className="form-label"
														htmlFor="form3Example1c"
													>
														Your Name
													</label>
												</div>
											</div>

											<div className="d-flex flex-row align-items-center mb-4">
												<i className="fas fa-envelope fa-lg me-3 fa-fw"></i>
												<div className="form-outline flex-fill mb-0">
													<input
														type="email"
														id="form3Example3c"
														className="form-control"
													/>
													<label
														className="form-label"
														htmlFor="form3Example3c"
													>
														Your Email
													</label>
												</div>
											</div>

											<div className="d-flex flex-row align-items-center mb-4">
												<i className="fas fa-lock fa-lg me-3 fa-fw"></i>
												<div className="form-outline flex-fill mb-0">
													<input
														type="password"
														id="form3Example4c"
														className="form-control"
													/>
													<label
														className="form-label"
														htmlFor="form3Example4c"
													>
														Password
													</label>
												</div>
											</div>

											<div className="dropdown">
												<a
													className="dropdown-toggle"
													role="button"
													id="dropdownMenuButton"
													data-mdb-toggle="dropdown"
													aria-expanded="false"
												>
													ROLE
												</a>
												<ul
													className="dropdown-menu"
													aria-labelledby="dropdownMenuButton"
												>
													<li>
														<a className="dropdown-item">Student</a>
													</li>
													<li>
														<a className="dropdown-item">Faculty</a>
													</li>
													<li>
														<a className="dropdown-item">Guest</a>
													</li>
												</ul>
											</div>

											{/* <div>
												<select
													className="dropdown"
													defaultValue={dropDownValue}
												>
													<option className="dropdown-menu" value="Orange">
														Orange
													</option>
													<option className="dropdown-menu" value="Radish">
														Radish
													</option>
													<option className="dropdown-menu" value="Cherry">
														Cherry
													</option>
												</select>
												<p>ROLE</p>
											</div> */}

											<div>
												<div className="form-check d-flex justify-content-center mb-5">
													Not a member? &nbsp;
													<Link to="/register" style={{ color: '#336EFF' }}>
														Register here!
													</Link>
												</div>

												<div className="d-flex justify-content-center mx-4 mb-3 mb-lg-4">
													<button
														type="button"
														className="btn btn-primary btn-lg"
													>
														Login
													</button>
												</div>
											</div>
										</form>
									</div>
									<div className="col-md-10 col-lg-6 col-xl-7 d-flex align-items-center order-1 order-lg-2">
										<img
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