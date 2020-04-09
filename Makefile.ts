import * as fs from "fs";
import yargs from "yargs";
import { Logger } from "tslog";
import * as execa from "execa";
import * as hasha from "hasha";
import * as readline from "readline";

// >>> CONFIGURATION
// -----------------------------------------------------------------------------
// Supply input to this task runner via CLI options or environment vars.
//
// see: https://yargs.js.org/
const config = yargs(process.argv.slice(2))
	.option("versionNo", { default: process.env["VERSION_NO"] ?? "0.0.0" })
	.option("buildDate", {
		default: process.env["DATE"] ?? new Date().toISOString(),
	})
	.option("commitUrl", {
		default:
			process.env["COMMIT_URL"] ??
			"https://github.com/owner/project/commit/hash",
	}).argv;

// >>> LOGGING
// -----------------------------------------------------------------------------
// All output from this task runner will be through this logger.
//
// see: https://tslog.js.org/
const logger = new Logger({
	displayInstanceName: false,
	displayLoggerName: false,
	displayFunctionName: false,
	displayFilePath: "hidden",
});

// >>> UTILS
// -----------------------------------------------------------------------------
// Functions that we use in our tasks.

/**
 * Executes a child process and outputs all stdio through the supplied logger.
 */
async function exe(
	log: Logger,
	args?: readonly string[],
	options?: execa.Options
) {
	const proc = execa(args[0], args.slice(1), options);
	readline
		.createInterface(proc.stdout)
		.on("line", (line: string) => log.info(line));
	readline
		.createInterface(proc.stderr)
		.on("line", (line: string) => log.warn(line));
	return await proc;
}

/**
 * Executes the `go build` command to produce a statically compiled binary.
 */
const goBuilder = async (os: "linux" | "windows", log: Logger = logger) =>
	await exe(
		log.getChildLogger({ prefix: ["go", "build", os] }),
		[
			"go",
			"build",
			"-v",
			"-ldflags",
			`-X main.versionNo=${config.versionNo} -X main.commitUrl=${config.commitUrl} -X main.buildDate=${config.buildDate}`,
			"-o",
			`./bin/dotfiles_${os}_amd64`,
		],
		{
			env: {
				CGO_ENABLED: "0",
				GOOS: os,
				GOARCH: "amd64",
			},
		}
	);

// >>> TASKS
// -----------------------------------------------------------------------------
export async function clean() {
	const log = logger.getChildLogger({ prefix: ["clean:"] });
	await fs.promises.rmdir("./bin", { recursive: true });
	log.info("rm -rf ./bin");
}

/**
 * With PowerShell call this like: `./make run '--' --version`
 * Bash of course understands with `--` means and doesn't require quotes.
 */
export async function run() {
	await build();

	const args = [
		`./bin/dotfiles_${
			process.platform === "win32" ? "windows" : "linux"
		}_amd64${process.platform === "win32" ? ".exe" : ""}`,
	];

	args.push(...(process.argv.slice(3) as string[]));

	await exe(logger.getChildLogger({ prefix: ["run:"] }), args);
}

export async function build() {
	await buildWindows();
}

export async function buildLinux() {
	const log = logger.getChildLogger({ prefix: ["buildLinux:"] });
	log.info("building go binary");
	await goBuilder("linux", log);
}

export async function buildWindows() {
	await buildLinux();

	const log = logger.getChildLogger({ prefix: ["buildWindows:"] });

	log.info(
		"copying ./bin/dotfiles_linux_amd64 => ./pkg/assets/files/dotfiles_linux_amd64"
	);
	await fs.promises.copyFile(
		"./bin/dotfiles_linux_amd64",
		"./pkg/assets/files/dotfiles_linux_amd64"
	);
	log.info(
		"copied ./bin/dotfiles_linux_amd64 => ./pkg/assets/files/dotfiles_linux_amd64"
	);

	try {
		log.info("building go binary");
		await goBuilder("windows", log);

		log.info("adding .exe suffix");
		await fs.promises.rename(
			"./bin/dotfiles_windows_amd64",
			"./bin/dotfiles_windows_amd64.exe"
		);
	} finally {
		log.info("deleting ./pkg/assets/files/dotfiles_linux_amd64");
		await fs.promises.rm("./pkg/assets/files/dotfiles_linux_amd64");
	}
}

export async function writeChecksums() {
	const log = logger.getChildLogger({ prefix: ["writeChecksums:"] });

	let checksumFile = "";

	for (let file of await fs.promises.readdir("./bin")) {
		const hash = await hasha.fromFile(`./bin/${file}`, {
			algorithm: "sha256",
		});
		checksumFile = `${checksumFile}${hash}  ${file}\n`;
	}

	await fs.promises.writeFile(
		"./bin/sha256_checksums.txt",
		checksumFile,
		"utf8"
	);

	log.info("written ./bin/sha256_checksums.txt");
}

export async function prepareRelease() {
	await clean();
	await build();
	await writeChecksums();
}

// >>> ENTRYPOINT
// -----------------------------------------------------------------------------
module.exports[config._[0]]
	.apply(null)
	.then(() => process.exit(0))
	.catch((e) => {
		logger.error(e);
		process.exit(1);
	});
