import React, { useEffect, useRef } from 'react';

const ParticleBackground = () => {
    const canvasRef = useRef(null);
    
    useEffect(() => {
        const canvas = canvasRef.current;
        const ctx = canvas.getContext('2d');
        let animationId;
        
        // Particle system implementation
        class Particle {
            constructor(x, y) {
                this.x = x;
                this.y = y;
                this.size = Math.random() * 2 + 0.5;
                this.speedX = Math.random() * 0.5 - 0.25;
                this.speedY = Math.random() * 0.5 - 0.25;
                this.opacity = Math.random() * 0.5 + 0.3;
            }
            
            update() {
                this.x += this.speedX;
                this.y += this.speedY;
                
                // Wrap around screen
                if (this.x > canvas.width) this.x = 0;
                if (this.x < 0) this.x = canvas.width;
                if (this.y > canvas.height) this.y = 0;
                if (this.y < 0) this.y = canvas.height;
            }
            
            draw() {
                ctx.save();
                ctx.globalAlpha = this.opacity;
                ctx.fillStyle = '#14b8a6';
                ctx.shadowBlur = 10;
                ctx.shadowColor = '#14b8a6';
                ctx.beginPath();
                ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
                ctx.fill();
                ctx.restore();
            }
        }
        
        // Create particles
        const particles = Array.from({ length: 100 }, () => 
            new Particle(
                Math.random() * canvas.width,
                Math.random() * canvas.height
            )
        );
        
        // Animation loop
        const animate = () => {
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            
            particles.forEach(particle => {
                particle.update();
                particle.draw();
            });
            
            // Draw connections
            particles.forEach((p1, i) => {
                particles.slice(i + 1).forEach(p2 => {
                    const distance = Math.hypot(p1.x - p2.x, p1.y - p2.y);
                    if (distance < 100) {
                        ctx.save();
                        ctx.globalAlpha = (1 - distance / 100) * 0.2;
                        ctx.strokeStyle = '#14b8a6';
                        ctx.lineWidth = 0.5;
                        ctx.beginPath();
                        ctx.moveTo(p1.x, p1.y);
                        ctx.lineTo(p2.x, p2.y);
                        ctx.stroke();
                        ctx.restore();
                    }
                });
            });
            
            animationId = requestAnimationFrame(animate);
        };
        
        // Handle resize
        const handleResize = () => {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
        };
        
        handleResize();
        window.addEventListener('resize', handleResize);
        animate();
        
        return () => {
            window.removeEventListener('resize', handleResize);
            cancelAnimationFrame(animationId);
        };
    }, []);
    
    return (
        <canvas
            ref={canvasRef}
            className="fixed inset-0 -z-10 opacity-50"
        />
    );
};

export default ParticleBackground;