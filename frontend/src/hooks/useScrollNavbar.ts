import { useState, useEffect, useRef } from "react";

export type NavbarState = "static" | "hidden" | "floating";

/**
 * Returns current state of the navbar depending on scroll direction.
 * Possible states:
 *  - "static"   — page is no the top (scrollY <= topThreshold)
 *  - "hidden"   — user scrolls down
 *  - "floating" — user scrolls up (not from the top)
 *
 * @param scrollThreshold  minimal scroll distance to trigger state change (defence from jitter)
 * @param topThreshold     distance from the top of the page to consider it as "static"
 */
export const useScrollNavbar = (
    scrollThreshold = 12,
    topThreshold = 8,
): NavbarState => {
    const [state, setState] = useState<NavbarState>("static");
    const lastScrollY = useRef(0);

    useEffect(() => {
        lastScrollY.current = window.scrollY;

        const handleScroll = () => {
            const currentY = window.scrollY;
            const diff = currentY - lastScrollY.current;

            if (currentY <= topThreshold) {
                setState("static");
            } else if (diff > scrollThreshold) {
                setState("hidden");
            } else if (diff < -scrollThreshold) {
                setState("floating");
            }

            lastScrollY.current = currentY;
        };

        window.addEventListener("scroll", handleScroll, { passive: true });
        return () => window.removeEventListener("scroll", handleScroll);
    }, [scrollThreshold, topThreshold]);

    return state;
};
