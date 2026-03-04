import { User } from "../api/userApi";
import "./css_files/top_nav.css";
import { useNavigate } from "react-router-dom";
interface NavBarProps {
  user: User | null;
}

const NavBar: React.FC<NavBarProps> = ({ user }) => {
  const navigate = useNavigate();
  return (
    <div className="navbar">
      <button className="sbu-button">
        <img
          src="/images/stonybrook.png"
          alt="Custom Icon"
          className="sbu-image"
          onClick={() => navigate("/home")}
        />
      </button>
      <div className="right">
      <button className="icon-button">
          <img
            src="/images/ticket.png"
            alt="Custom Icon"
            className="icon-image"
            onClick={() => navigate("/tickets")}
          />
        </button>

        {user?.Role === "admin" && (
          <button className="icon-button">
            <img
              src="/images/edit.png"
              alt="Custom Icon"
              className="icon-image"
              onClick={() => navigate("/edit")}
            />
          </button>
        )}
        <button className="icon-button">
          <img
            src="/images/profile_icon.png"
            alt="Custom Icon"
            className="icon-image"
            onClick={() => navigate("/profile")}
          />
        </button>
      </div>
    </div>
  );
};

export default NavBar;
