import pg from "pg";
import { DATABASE_URL } from "../env.js";

const pool = new pg.Pool({ connectionString: DATABASE_URL.value() });

export { pool };
