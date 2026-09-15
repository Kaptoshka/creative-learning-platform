import React, { useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { BookOpen, Home, LogIn, LogOut } from "lucide-react";
import { AuthContext } from "@/context/AuthContext";
import { useScrollNavbar } from "@/hooks/useScrollNavbar";
import MobileMenuToggleButton from "./MobileMenuToggleButton";
import NavigationLink from "./NavigationLink";
import Button from "@/components/Button";
import styles from "./Navigation.module.scss";

const Navigation = () => {
    const { user, logout } = useContext(AuthContext);
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
    const navigate = useNavigate();
    const navState = useScrollNavbar();

    const handleLogout = () => {
        logout();
        navigate("/", { viewTransition: true });
    };

    return (
        <nav className={`${styles.navigation} ${styles[navState]}`}>
            <div className={styles.navigationContainer}>
                <div className={styles.navigationContent}>
                    <div className={styles.navigationLogo}>
                        <BookOpen className={styles.navigationLogoIcon} />
                        <span className={styles.navigationLogoText}>
                            CreativeLearning
                        </span>
                    </div>
                    {!isMobileMenuOpen && (
                        <div className={styles.navigationMenu}>
                            <NavigationLink to="/" end icon={Home}>
                                Главная
                            </NavigationLink>

                            {user ? (
                                <>
                                    <NavigationLink
                                        to="/dashboard"
                                        icon={BookOpen}
                                    >
                                        Личный кабинет
                                    </NavigationLink>
                                    <Button
                                        onClick={handleLogout}
                                        variant="navigation"
                                        icon={LogOut}
                                    >
                                        {" "}
                                        Выйти
                                    </Button>
                                </>
                            ) : (
                                <NavigationLink to="/auth" icon={LogIn}>
                                    Войти
                                </NavigationLink>
                            )}
                        </div>
                    )}

                    <div className={styles.navigationMobileToggle}>
                        <MobileMenuToggleButton
                            isOpen={isMobileMenuOpen}
                            onClick={() =>
                                setIsMobileMenuOpen(!isMobileMenuOpen)
                            }
                        />
                    </div>
                </div>

                {isMobileMenuOpen && (
                    <div className={styles.navigationMobileMenu}>
                        <NavigationLink
                            to="/"
                            end
                            onClick={() => setIsMobileMenuOpen(false)}
                        >
                            Главная
                        </NavigationLink>
                        {user ? (
                            <>
                                <NavigationLink
                                    to="/dashboard"
                                    onClick={() => setIsMobileMenuOpen(false)}
                                >
                                    Личный кабинет
                                </NavigationLink>
                                <Button
                                    onClick={() => {
                                        handleLogout();
                                        setIsMobileMenuOpen(false);
                                    }}
                                    variant="mobile-navigation"
                                >
                                    Выйти
                                </Button>
                            </>
                        ) : (
                            <NavigationLink
                                to="/auth"
                                onClick={() => setIsMobileMenuOpen(false)}
                            >
                                Войти
                            </NavigationLink>
                        )}
                    </div>
                )}
            </div>
        </nav>
    );
};

export default Navigation;
