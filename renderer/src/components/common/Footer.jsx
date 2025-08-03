import { motion } from 'framer-motion';
import footerImg from '../../assets/images/footer-logo.jpg';


const Footer = () => {
    return (
        <motion.footer
            className="fixed bottom-0 left-0 w-full glass border-t border-teal-500/20 z-50"
            style={{ WebkitAppRegion: 'drag' }}
            initial={{ y: 100 }}
            animate={{ y: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
        >
            <div className="container mx-auto px-4 py-3 flex flex-col items-center justify-center">
                <motion.div
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    transition={{ duration: 0.5, delay: 0.5 }}
                    className="relative"
                >
                    <div className="absolute inset-0 blur-xl bg-teal-500/30 rounded-full" />
                    <img
                        src={footerImg}
                        alt="Logo"
                        className="h-8 w-8 mb-2 rounded-full border border-teal-500/50 relative z-10"
                    />
                </motion.div>
                <motion.span 
                    className="text-sm text-gray-400"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.5, delay: 0.7 }}
                >
                    &copy; {new Date().getFullYear()} CanIFly. All rights reserved.
                </motion.span>
            </div>
        </motion.footer>
    );
};

export default Footer;
