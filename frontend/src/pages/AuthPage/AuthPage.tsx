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

    return (
        <div className={styles.authPage}>
            {/* Left decorative panel */}
            <div className={styles.authPanel}>
                <div className={styles.authPanelContent}>
                    <div className={styles.authPanelLogo}>
                        <svg
                            width="32"
                            height="32"
                            viewBox="0 0 32 32"
                            fill="none"
                        >
                            <rect
                                width="14"
                                height="14"
                                rx="3"
                                fill="white"
                                fillOpacity="0.9"
                            />
                            <rect
                                x="18"
                                width="14"
                                height="14"
                                rx="3"
                                fill="white"
                                fillOpacity="0.6"
                            />
                            <rect
                                y="18"
                                width="14"
                                height="14"
                                rx="3"
                                fill="white"
                                fillOpacity="0.6"
                            />
                            <rect
                                x="18"
                                y="18"
                                width="14"
                                height="14"
                                rx="3"
                                fill="white"
                                fillOpacity="0.9"
                            />
                        </svg>
                    </div>
                    <div className={styles.authPanelText}>
                        <h1>
                            Образование,
                            <br />
                            которое работает
                        </h1>
                        <p>
                            Управляйте заданиями, отслеживайте прогресс и
                            получайте обратную связь в одном месте.
                        </p>
                    </div>
                    <div className={styles.authPanelDecor}>
                        <div className={styles.decor1} />
                        <div className={styles.decor2} />
                        <div className={styles.decor3} />
                        <div className={styles.decor4} />
                        <div className={styles.decor5} />
                    </div>
                </div>
            </div>

            {/* Right form panel */}
            <div className={styles.authFormPanel}>
                <div className={styles.authFormWrapper}>
                    <div className={styles.authFormHeader}>
                        <h2>
                            {isLoginMode
                                ? "Войти в аккаунт"
                                : "Создать аккаунт"}
                        </h2>
                        <p>
                            {isLoginMode
                                ? "Добро пожаловать обратно!"
                                : "Присоединяйтесь к нам!"}
                        </p>
                    </div>

                    {error && <div className={styles.authError}>{error}</div>}

                    <form
                        onSubmit={handleSubmit}
                        noValidate
                        className={styles.authForm}
                    >
                        {!isLoginMode && (
                            <>
                                <div className={styles.authFieldRow}>
                                    {/*
                                     * ВАЖНО: input идёт ДО label.
                                     * Floating label работает через CSS-селектор input:focus ~ label
                                     * и input:not(:placeholder-shown) ~ label
                                     */}
                                    <div className={styles.authField}>
                                        <input
                                            id="firstName"
                                            type="text"
                                            placeholder=" "
                                            value={firstName}
                                            onChange={(e) =>
                                                setFirstName(e.target.value)
                                            }
                                            required
                                            disabled={isLoading}
                                        />
                                        <label htmlFor="firstName">Имя</label>
                                    </div>
                                    <div className={styles.authField}>
                                        <input
                                            id="lastName"
                                            type="text"
                                            placeholder=" "
                                            value={lastName}
                                            onChange={(e) =>
                                                setLastName(e.target.value)
                                            }
                                            required
                                            disabled={isLoading}
                                        />
                                        <label htmlFor="lastName">
                                            Фамилия
                                        </label>
                                    </div>
                                </div>
                                <div className={styles.authField}>
                                    <input
                                        id="middleName"
                                        type="text"
                                        placeholder=" "
                                        value={middleName}
                                        onChange={(e) =>
                                            setMiddleName(e.target.value)
                                        }
                                        disabled={isLoading}
                                    />
                                    <label htmlFor="middleName">Отчество</label>
                                </div>
                            </>
                        )}

                        <div className={styles.authField}>
                            <input
                                id="email"
                                type="email"
                                placeholder=" "
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                                disabled={isLoading}
                            />
                            <label htmlFor="email">Email</label>
                        </div>

                        <div className={styles.authField}>
                            <input
                                id="password"
                                type="password"
                                placeholder=" "
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                                disabled={isLoading}
                            />
                            <label htmlFor="password">Пароль</label>
                        </div>

                        {!isLoginMode && (
                            <div className={styles.authField}>
                                <input
                                    id="confirmPassword"
                                    type="password"
                                    placeholder=" "
                                    value={confirmPassword}
                                    onChange={(e) =>
                                        setConfirmPassword(e.target.value)
                                    }
                                    required
                                    disabled={isLoading}
                                />
                                <label htmlFor="confirmPassword">
                                    Подтвердите пароль
                                </label>
                            </div>
                        )}

                        <Button
                            type="submit"
                            variant="primary"
                            fullWidth
                            disabled={isLoading}
                            className={styles.authSubmit}
                        >
                            {isLoading
                                ? "Загрузка..."
                                : isLoginMode
                                  ? "Войти"
                                  : "Зарегистрироваться"}
                        </Button>
                    </form>

                    <div className={styles.authSwitch}>
                        <span>
                            {isLoginMode
                                ? "Нет аккаунта? "
                                : "Уже есть аккаунт? "}
                        </span>
                        <button onClick={handleModeSwitch} disabled={isLoading}>
                            {isLoginMode ? "Зарегистрируйтесь" : "Войдите"}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AuthPage;
