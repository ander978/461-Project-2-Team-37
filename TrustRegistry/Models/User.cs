using System.Collections.Generic;

namespace TrustRegistry.Models
{
    // Model for a user
    public class User
    {
        // The name of the user
        public string Name { get; set; }
        
        // Whether or not the user is an admin
        public bool IsAdmin { get; set; }

        // Dictionary to hold user authentication information
        private Dictionary<string, string> _UserAuthInf;

        // Getter for user authentication information
        public Dictionary<string, string> GetUserAuthInf()
        {
            return _UserAuthInf;
        }

        // Setter for user authentication information
        public void SetUserAuthInf(Dictionary<string, string> UserAuthInf)
        {
            _UserAuthInf = UserAuthInf;
        }

        // Constructor for creating a new user object
        public User(string name, bool isAdmin)
        {
            Name = name;
            IsAdmin = isAdmin;
        }
    }
}