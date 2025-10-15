<script lang="ts">
	import { onMount } from 'svelte';

	let content = $state<HTMLDivElement>();
	let container = $state<HTMLDivElement>();

	let scrollRatio = $state(0);
	let thumbHeight = $state(0);

	function updateThumb() {
		const { scrollTop, scrollHeight, clientHeight } = content!;
		scrollRatio = scrollTop / (scrollHeight - clientHeight);
		thumbHeight = (clientHeight / scrollHeight) * container!.clientHeight;
	}

	function scrollBy(delta: number) {
		content!.scrollTop += delta;
		updateThumb();
	}

	function handleThumbDrag(e: MouseEvent) {
		e.preventDefault();
		const startY = e.clientY;
		const startScroll = content!.scrollTop;

		function move(ev: MouseEvent) {
			const diff = ev.clientY - startY;
			const scrollDiff = (diff / container!.clientHeight) * content!.scrollHeight;
			content!.scrollTop = startScroll + scrollDiff;
			updateThumb();
		}

		function stop() {
			window.removeEventListener('mousemove', move);
			window.removeEventListener('mouseup', stop);
		}

		window.addEventListener('mousemove', move);
		window.addEventListener('mouseup', stop);
	}

	onMount(() => {
		updateThumb();
	});

	let { children } = $props();
</script>

<div class="relative flex h-full rounded-lg bg-custom-dark-white p-4">
	<!-- Scrollable content -->
	<div
		bind:this={content}
		aria-owns="main"
		class="no-scrollbar flex-1 overflow-y-scroll"
		onscroll={updateThumb}
	>
		{@render children?.()}
	</div>

	<div
		hidden={content != undefined ? content.scrollHeight <= content.clientHeight : false}
		bind:this={container}
		class="my-4 ml-6 flex w-[1px] flex-col items-center justify-between rounded-lg bg-[#7B7B7B]"
	>
		<button
			aria-label="scroll up"
			onclick={() => scrollBy(-50)}
			class="-mt-4 text-gray-500 hover:text-black"
		>
			<svg width="14" height="8" viewBox="0 0 17 10" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path
					fill-rule="evenodd"
					clip-rule="evenodd"
					d="M7.21748 0.357194C7.44543 0.128483 7.75455 -3.68445e-07 8.07686 -3.54357e-07C8.39918 -3.40268e-07 8.7083 0.128483 8.93624 0.357194L15.8125 7.25864C15.9286 7.37118 16.0212 7.50579 16.0849 7.65464C16.1486 7.80348 16.1821 7.96357 16.1835 8.12556C16.1849 8.28754 16.1542 8.44819 16.0931 8.59812C16.032 8.74805 15.9417 8.88427 15.8276 8.99881C15.7134 9.11336 15.5777 9.20395 15.4283 9.26529C15.279 9.32663 15.1189 9.3575 14.9575 9.35609C14.7961 9.35469 14.6366 9.32103 14.4883 9.25709C14.34 9.19316 14.2059 9.10021 14.0937 8.98369L8.07686 2.94478L2.05998 8.98369C1.83072 9.20592 1.52368 9.32889 1.20497 9.32611C0.886262 9.32333 0.581392 9.19503 0.356022 8.96883C0.130653 8.74264 0.00281488 8.43665 4.54216e-05 8.11677C-0.00272403 7.7969 0.119795 7.48873 0.341214 7.25864L7.21748 0.357194Z"
					fill="#7B7B7B"
				/>
			</svg>
		</button>

		<div class="relative w-full flex-1">
			<div
				role="scrollbar"
				aria-controls="main"
				aria-valuenow="0"
				tabindex={0}
				class="absolute left-1/2 -translate-x-1/2 transform rounded-full bg-gray-500"
				style="height: {thumbHeight}px; top: calc({scrollRatio} * (100% - {thumbHeight}px)); width: 6px;"
				onmousedown={handleThumbDrag}
			></div>
		</div>

		<button
			aria-label="scroll down"
			onclick={() => scrollBy(50)}
			class="-mb-4 text-gray-500 hover:text-black"
		>
			<svg width="14" height="8" viewBox="0 0 17 10" fill="none" xmlns="http://www.w3.org/2000/svg">
				<path
					fill-rule="evenodd"
					clip-rule="evenodd"
					d="M7.21748 8.99827C7.44543 9.22699 7.75455 9.35547 8.07686 9.35547C8.39918 9.35547 8.7083 9.22699 8.93624 8.99827L15.8125 2.09683C15.9286 1.98429 16.0212 1.84967 16.0849 1.70083C16.1486 1.55199 16.1821 1.3919 16.1835 1.22991C16.1849 1.06792 16.1542 0.907279 16.0931 0.757347C16.032 0.607416 15.9417 0.471201 15.8276 0.356654C15.7134 0.242106 15.5777 0.151519 15.4283 0.0901776C15.279 0.0288353 15.1189 -0.00203223 14.9575 -0.000624603C14.7961 0.000783027 14.6366 0.0344372 14.4883 0.0983754C14.34 0.162314 14.2059 0.255255 14.0937 0.371776L8.07686 6.41069L2.05998 0.371776C1.83072 0.149546 1.52368 0.0265795 1.20497 0.0293586C0.886262 0.0321385 0.581392 0.160443 0.356022 0.386637C0.130653 0.612833 0.00281488 0.918819 4.54216e-05 1.23869C-0.00272403 1.55857 0.119795 1.86674 0.341214 2.09683L7.21748 8.99827Z"
					fill="#7B7B7B"
				/>
			</svg>
		</button>
	</div>
</div>

<style>
	.no-scrollbar::-webkit-scrollbar {
		display: none;
	}
	.no-scrollbar {
		-ms-overflow-style: none;
		scrollbar-width: none;
	}
</style>
