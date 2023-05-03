namespace TrustRegistry.Models 
{
    public class AuthRequest 
    {
        public User User { get; set; } // User object containing user information
        public UserAuthInf UserAuthInf { get; set; } // UserAuthenticationInfo object containing user authentication information
    }
}