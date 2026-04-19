import React, { useState, useContext } from "react";
import { flushSync } from "react-dom";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "@/context/AuthContext";
import { BookOpen } from "lucide-react";
import axios from "axios";
import Button from "@/components/Button";
import styles from "./AuthPage.module.scss";
import { config } from "@/config";

const AuthPage = () => {
    const [mode, setMode] = useState("login");

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [middleName, setMiddleName] = useState("");

    const [error, setError] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const { login } = useContext(AuthContext);
    const navigate = useNavigate();

    const handleModeSwitch = () => {
        setEmail("");
        setPassword("");
        setConfirmPassword("");
        setFirstName("");
        setLastName("");
        setMiddleName("");
        setError("");

        const toggleMode = () => {
            setMode((prev) => (prev === "login" ? "register" : "login"));
        };

        if (!document.startViewTransition) {
            toggleMode();
            return;
        }

        document.startViewTransition(() => {
            flushSync(() => {
                toggleMode();
            });
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setError("");

        if (mode == "register" && password !== confirmPassword) {
            setError("Пароли не совпадают.");
            setIsLoading(false);
            return;
        }

        try {
            let response;
            if (mode === "login") {
                response = await axios.post(`${config.ssoAPIURL}/auth/login`, {
                    email: email,
                    password: password,
                    app_id: +config.appId,
                });
            } else {
                response = await axios.post(
                    `${config.ssoAPIURL}/auth/register`,
                    {
                        email: email,
                        password: password,
                        first_name: firstName,
                        last_name: lastName,
                        middle_name: middleName,
                    },
                );
            }

            if (response.data) {
                if (response.data.token) {
                    login(response.data.token);
                } else if (response.data.user_id) {
                    navigate("/auth", { viewTransition: true });
                }
                navigate("/tasks", { viewTransition: true });
            } else {
                setError("login failed: no token was received from the server");
            }
        } catch (err) {
            if (err.response && err.response.data && err.response.data.error) {
                setError(err.response.data.error);
            } else {
                setError(
                    "login failed: please check your credentials or network connection",
                );
            }
            console.error("login failed: ", err);
        } finally {
            setIsLoading(false);
        }
    };

    const isLoginMode = mode === "login";
    const title = isLoginMode ? "Логин" : "Создать аккаунт";
    const buttonText = isLoginMode ? "Логин" : "Регистрация";
    const switchButtonText = isLoginMode
        ? "Нет аккаунта? Зарегистрируйтесь"
        : "Уже есть аккаунт? Войдите";

    return (
        <div className={styles.authPage}>
            <form onSubmit={handleSubmit} noValidate>
                <div className={styles.authForm}>
                    <div className={styles.authPageHeader}>
                        <BookOpen className={styles.authPageLogo} />
                        <h2 className={styles.authPageTitle}>{title}</h2>
                        <p className={styles.authPageSubtitle}>
                            {isLoginMode
                                ? "Добро пожаловать обратно!"
                                : "Присоединяйтесь к нам!"}
                        </p>
                    </div>
                    {error && <p className="error-message">{error}</p>}
                    <div className={styles.authPageFormContainer}>
                        <div className={styles.authFormContent}>
                            <div className={styles.authFormColumns}>
                                {!isLoginMode && (
                                    <div className={styles.authFormColumn}>
                                        <div className="form-group">
                                            <label
                                                htmlFor="name"
                                                className="form-group__label"
                                            >
                                                Имя
                                            </label>
                                            <input
                                                id="name"
                                                type="text"
                                                className="form-group__input"
                                                placeholder="Введите имя"
                                                value={firstName}
                                                onChange={(e) =>
                                                    setFirstName(e.target.value)
                                                }
                                                required
                                            />
                                        </div>
                                        <div className="form-group">
                                            <label
                                                htmlFor="name"
                                                className="form-group__label"
                                            >
                                                Фамилия
                                            </label>
                                            <input
                                                id="name"
                                                type="text"
                                                className="form-group__input"
                                                placeholder="Введите фамилию"
                                                value={lastName}
                                                onChange={(e) =>
                                                    setLastName(e.target.value)
                                                }
                                                required
                                            />
                                        </div>
                                        <div className="form-group">
                                            <label
                                                htmlFor="name"
                                                className="form-group__label"
                                            >
                                                Отчество
                                            </label>
                                            <input
                                                id="middleName"
                                                type="text"
                                                className="form-group__input"
                                                placeholder="Введите отчество"
                                                value={middleName}
                                                onChange={(e) =>
                                                    setMiddleName(
                                                        e.target.value,
                                                    )
                                                }
                                                required
                                                disabled={isLoading}
                                            />
                                        </div>
                                    </div>
                                )}
                                <div className={styles.authFormColumn}>
                                    <div className="form-group">
                                        <label
                                            htmlFor="email"
                                            className="form-group__label"
                                        >
                                            Email
                                        </label>
                                        <input
                                            id="email"
                                            type="email"
                                            className="form-group__input"
                                            placeholder="Введите email"
                                            value={email}
                                            onChange={(e) =>
                                                setEmail(e.target.value)
                                            }
                                            required
                                            disabled={isLoading}
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label
                                            htmlFor="password"
                                            className="form-group__label"
                                        >
                                            Пароль
                                        </label>
                                        <input
                                            id="password"
                                            type="password"
                                            className="form-group__input"
                                            placeholder="Введите пароль"
                                            value={password}
                                            onChange={(e) =>
                                                setPassword(e.target.value)
                                            }
                                            required
                                            disabled={isLoading}
                                        />
                                    </div>
                                    {!isLoginMode && (
                                        <div className="form-group">
                                            <label
                                                htmlFor="confirmPassword"
                                                className="form-group__label"
                                            >
                                                Подтвердите пароль
                                            </label>
                                            <input
                                                id="confirmPassword"
                                                type="password"
                                                className="form-group__input"
                                                placeholder="Подтвердите пароль"
                                                value={confirmPassword}
                                                onChange={(e) =>
                                                    setConfirmPassword(
                                                        e.target.value,
                                                    )
                                                }
                                                required
                                                disabled={isLoading}
                                            />
                                        </div>
                                    )}
                                </div>
                            </div>

                            <Button
                                type="submit"
                                variant="primary"
                                fullWidth
                                disabled={isLoading}
                            >
                                {isLoading ? "Загрузка..." : buttonText}
                            </Button>
                            <div className={styles.authFormFooter}>
                                <Button
                                    type="button"
                                    className={styles.authFormSwitch}
                                    onClick={handleModeSwitch}
                                    disabled={isLoading}
                                >
                                    {switchButtonText}
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>
            </form>
        </div>
    );
};

export default AuthPage;
