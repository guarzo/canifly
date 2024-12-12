import footerImg from '../../assets/images/footer-logo.jpg';


const Footer = () => {
    return (
        <footer
            className="fixed bottom-0 left-0 w-full bg-gradient-to-r from-gray-900 to-gray-800 text-teal-200 py-4 shadow-inner border-t-4 border-teal-500 z-50"
            style={{ WebkitAppRegion: 'drag' }}
        >
            <div className="container mx-auto px-4 flex flex-col items-center justify-center">
                <img
                    src={footerImg}
                    alt="Logo"
                    className="h-8 w-8 mb-2 rounded-full border-2 border-teal-500"
                />
                <span className="text-sm">
                    &copy; {new Date().getFullYear()} Can I Fly? All rights reserved.
                </span>
            </div>
        </footer>
    );
};

export default Footer;
